package options

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var ProgramOptions Options
var Version string = "0.1"

func init() {
	// Connection parameters.
	flag.StringVar(&ProgramOptions.User, "user", "root", "user to access database.")
	flag.StringVar(&ProgramOptions.User, "u", "root", "user to access database.")

	flag.StringVar(&ProgramOptions.Pass, "pass", "", "password to access database.")
	flag.StringVar(&ProgramOptions.Pass, "p", "", "password to access database.")

	flag.StringVar(&ProgramOptions.Host, "host", "localhost:3306", "host[:port]")
	flag.StringVar(&ProgramOptions.Host, "h", "localhost:3306", "host[:port]")

	flag.StringVar(&ProgramOptions.SchemaName, "database", "", "database name")
	flag.StringVar(&ProgramOptions.SchemaName, "d", "", "database name")

	// Modes
	flag.BoolVar(&ProgramOptions.Version, "version", false, "display version info and exit.")
	flag.StringVar(&ProgramOptions.DumpFile, "dump", "", "xml dump file of a database structure.")
	flag.StringVar(&ProgramOptions.StructureFile, "restore", "", "xml dump file of a database structure.")

	// Verbosity level
	flag.BoolVar(&ProgramOptions.VeryQuiet, "qq", false, "don't produce any output.")
	flag.BoolVar(&ProgramOptions.Quiet, "q", false, "only show errors.")
	flag.BoolVar(&ProgramOptions.Verbose, "v", false, "verbose output.")
	flag.BoolVar(&ProgramOptions.VeryVerbose, "vv", false, "very verbose output.")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [inputfile]\n", os.Args[0])
	flag.PrintDefaults()
}

func Parse() bool {
	flag.Usage = usage
	flag.Parse()

	if ! strings.ContainsRune(ProgramOptions.Host, ':') {
		ProgramOptions.Host += ":3306"
	}
	// Activate 'debug' mode.
	ProgramOptions.Debug = ProgramOptions.VeryVerbose

	// Can't give both flags at the same time.
	if "" != ProgramOptions.DumpFile && "" != ProgramOptions.StructureFile {
		fmt.Fprintln(os.Stderr, "Can't do --dump and --restore at the same time!\n")
		usage()
		return false
	}

	// Need to give one of them.
	if "" == ProgramOptions.DumpFile && "" == ProgramOptions.StructureFile && ! ProgramOptions.Version {
		fmt.Fprintln(os.Stderr, "You didn't tell me what to do!\n")
		usage()
		return false
	}

	// Need to give database name.
	if "" == ProgramOptions.SchemaName {
		fmt.Fprintln(os.Stderr, "Missing database name.\n")
		usage()
		return false
	}

	return true
}
