package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
)

type (
	// Restore includes logic to restore a db from s3 or your home directory
	Restore struct {
		bucket       string
		cmd          *cobra.Command
		extract      int
		prefix       string
		remote       bool
		s3FilePrefix string
	}
)

func generateRestoreCommand() *cobra.Command {
	return Restore{}.command()
}

func (r Restore) command() *cobra.Command {
	if r.cmd == nil {
		cmd := &cobra.Command{
			Use:   "restore",
			Short: "Restore from an org extract",
			Args:  cobra.MaximumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				var extractArg int
				var err error
				if len(args) == 1 {
					extractArg, err = strconv.Atoi(args[0])
				}
				if extractArg == 0 || err != nil {
					extractArg = 4
				}
				r.extract = extractArg
				r.dropConnections()
				r.dropDb()
				r.createDb()

				if r.remote {
					err := r.downloadDump()
					onError(err)
				}
				r.restoreDb()
			},
		}

		cmd.Flags().StringVarP(&r.bucket, "bucket", "b", "hqdatabase-tracker-production", "s3 bucket that contains org extracts")
		cmd.Flags().StringVarP(&r.prefix, "prefix", "p", "agencieshq", "prefix of sql dump")
		cmd.Flags().BoolVarP(&r.remote, "remote", "r", false, "restore extract from AWS")
		cmd.Flags().StringVarP(&r.s3FilePrefix, "s3", "s", "agencieshq", "s3 prefix of sql dump")

		r.cmd = cmd
	}
	return r.cmd
}

func (r Restore) database() string {
	extractNo := fmt.Sprintf("%d", r.extract)
	return strings.Join([]string{r.prefix, extractNo}, "_")
}

func (r Restore) dump() string {
	return r.database() + ".dump"
}

func (r Restore) dropConnections() {
	sql := fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '%s'
			AND pid <> pg_backend_pid();
	`, r.database())

	cmd := OSCommand{
		Cmd:         "psql",
		Args:        []string{"-c", sql},
		Description: "drop connections",
	}
	cmd.Run()
}

func (r Restore) dropDb() {
	database := r.database()
	cmd := OSCommand{
		Cmd:         "dropdb",
		Args:        []string{database},
		Description: "drop db",
		SilentError: true,
	}
	cmd.Run()
}

func (r Restore) createDb() {
	database := r.database()
	cmd := OSCommand{
		Cmd:         "createdb",
		Args:        []string{database},
		Description: "create db",
	}
	cmd.Run()
}

func (r Restore) restoreDb() {
	database := r.database()
	dump := getHome() + "/" + r.dump()
	cmd := OSCommand{
		Cmd:         "pg_restore",
		Args:        []string{"-O", "-j", "8", "-d", database, dump},
		Description: "restore db",
	}
	cmd.Run()
}

func (r Restore) downloadDump() error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(RootConfig.AWSKey(), RootConfig.AWSSecretKey(), ""),
	})

	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	filename := r.s3FilePrefix + "_" + fmt.Sprintf("%d", r.extract) + ".dump"
	f, err := os.Create(getHome() + "/" + filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %v", filename, err)
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}
	fmt.Printf("file downloaded, %d bytes\n", n)
	return err
}
