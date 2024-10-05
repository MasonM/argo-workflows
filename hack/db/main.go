package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/upper/db/v4"
	mysqladp "github.com/upper/db/v4/adapter/mysql"
	postgresqladp "github.com/upper/db/v4/adapter/postgresql"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/uuid"

	"github.com/argoproj/argo-workflows/v3/persist/sqldb"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v3/util/instanceid"
)

var session db.Session

func main() {
	var dbtype, dsn string
	rootCmd := &cobra.Command{
		Use:   "db",
		Short: "CLI for developers to use when working on the DB locally",
	}
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		session, err = createDBSession(dbtype, dsn)
		return
	}
	rootCmd.PersistentFlags().StringVarP(&dbtype, "driver", "d", "postgresql", "Database type (mysql or postgresql)")
	rootCmd.PersistentFlags().StringVarP(&dsn, "dsn", "c", "postgres://postgres@localhost:5432/postgres", "DSN connection string")
	rootCmd.AddCommand(NewMigrateCommand())
	rootCmd.AddCommand(NewFakeDataCommand())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewMigrateCommand() *cobra.Command {
	var cluster, table string
	migrationCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Force DB migration for given cluster/table",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sqldb.NewMigrate(session, cluster, table).Exec(context.Background())
		},
	}
	migrationCmd.Flags().StringVar(&cluster, "cluster", "default", "Cluster name")
	migrationCmd.Flags().StringVar(&table, "table", "argo_workflows", "Table name")
	return migrationCmd
}

func NewFakeDataCommand() *cobra.Command {
	var seed, rows, numClusters, numNamespaces int
	fakeDataCmd := &cobra.Command{
		Use:   "fake-archived-workflows",
		Short: "Insert randomly-generated workflows into argo_archived_workflows, for testing purposes",
		RunE: func(cmd *cobra.Command, args []string) error {
			rand.Seed(int64(seed))
			clusters := randomStringArray(numClusters)
			namespaces := randomStringArray(numNamespaces)
			fmt.Printf("Using seed %d\nClusters: %v\nNamespaces: %v\n", seed, clusters, namespaces)

			instanceIDService := instanceid.NewService("")

			for i := 0; i < rows; i++ {
				wf := randomWorkflow(namespaces)
				cluster := clusters[rand.Intn(len(clusters))]
				wfArchive := sqldb.NewWorkflowArchive(session, cluster, "", instanceIDService)
				if err := wfArchive.ArchiveWorkflow(wf); err != nil {
					return err
				}
			}
			fmt.Printf("Inserted %d rows\n", rows)
			return nil
		},
	}
	fakeDataCmd.Flags().IntVar(&seed, "seed", rand.Int(), "Random number seed")
	fakeDataCmd.Flags().IntVar(&rows, "rows", 10, "Number of rows to insert")
	fakeDataCmd.Flags().IntVar(&numClusters, "clusters", 1, "Number of cluster names to autogenerate")
	fakeDataCmd.Flags().IntVar(&numNamespaces, "namespaces", 5, "Number of namespaces to autogenerate")
	return fakeDataCmd
}

func createDBSession(dbtype, dsn string) (db.Session, error) {
	if dbtype == "postgresql" {
		url, err := postgresqladp.ParseURL(dsn)
		if err != nil {
			return nil, err
		}
		return postgresqladp.Open(url)
	} else {
		url, err := mysqladp.ParseURL(dsn)
		if err != nil {
			return nil, err
		}
		return mysqladp.Open(url)
	}
}

func randomStringArray(length int) []string {
	var result []string
	for i := 0; i < length; i++ {
		result = append(result, rand.String(rand.IntnRange(5, 20)))
	}
	return result
}

const wfTmpl = `apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  name: %s
  namespace: %s
  uid: %s
  labels:
    workflows.argoproj.io/completed: "true"
    workflows.argoproj.io/phase: %s
spec:
  templates:
    - name: main
      container:
        image: argoproj/argosay:v2
        args: ["echo", "a"]
status:
  startedAt: 2000-01-01T00:00:00Z
  finishedAt: 2000-01-02T00:00:00Z
  phase: %[4]s
`

func randomPhase() wfv1.WorkflowPhase {
	phases := []wfv1.WorkflowPhase{
		wfv1.WorkflowSucceeded,
		wfv1.WorkflowFailed,
		wfv1.WorkflowError,
	}
	return phases[rand.Intn(len(phases))]
}

func randomWorkflow(namespaces []string) *wfv1.Workflow {
	wfString := fmt.Sprintf(wfTmpl,
		rand.String(rand.IntnRange(10, 30)),
		namespaces[rand.Intn(len(namespaces))],
		uuid.NewUUID(),
		randomPhase(),
	)
	return wfv1.MustUnmarshalWorkflow(wfString)
}
