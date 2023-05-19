package neo4j_tracing

const tracerName = "github.com/raito-io/neo4j_tracing"
const serviceID = "neo4j"

func spanName(operation string) string {
	return serviceID + "." + operation
}
