package main

import (
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

const (
	numClusters   = 5
	numNamespaces = 5
)

func main() {
	var dbtype, dsn, table string
	var rows int

	command := &cobra.Command{
		Use:   "db-data-generator",
		Short: "CLI to generate fake/test data and insert it into the database",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			session, err := createDBSession(dbtype, dsn)
			if err != nil {
				return
			}

			clusters := randomArray(numClusters)
			namespaces := randomArray(numNamespaces)
			instanceIDService := instanceid.NewService("")

			for i := 0; i < rows; i++ {
				wf := randomWorkflow(namespaces)
				cluster := clusters[rand.Intn(len(clusters))]
				wfArchive := sqldb.NewWorkflowArchive(session, cluster, "", instanceIDService)
				if err := wfArchive.ArchiveWorkflow(wf); err != nil {
					return err
				}
			}
			fmt.Printf("Rows: %v, Error: %v\n", rows, err)
			return
		},
	}
	command.Flags().StringVarP(&dbtype, "driver", "d", "postgresql", "Database type (mysql or postgresql)")
	command.Flags().StringVarP(&dsn, "dsn", "c", "postgres://postgres@localhost:5432/postgres", "DSN connection string")
	command.Flags().StringVarP(&table, "table", "t", "argo_archived_workflows", "Table to populate")
	command.Flags().IntVarP(&rows, "rows", "r", 10, "Number of rows to insert")
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func randomArray(length int) []string {
	var result []string
	for i := 0; i < length; i++ {
		result = append(result, rand.String(rand.IntnRange(5, 20)))
	}
	return result
}

var wfTmpl = `apiVersion: argoproj.io/v1alpha1
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
