package behavior

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/senrok/yadal"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/layers"
	"github.com/senrok/yadal/providers/fs"
	"github.com/senrok/yadal/providers/s3"
	"go.uber.org/zap"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type testFunc func(t *testing.T, op *yadal.Operator)

type strategy func(op *yadal.Operator) bool

type testSet struct {
	name string
	strategy
	tests []testFunc
}

var (
	providers = []string{"s3", "fs"}
	tests     = []testSet{
		{
			name: "basic",
			strategy: func(op *yadal.Operator) bool {
				return true
			},
			tests: baseTests,
		},
		{
			name: "readWrite",
			strategy: func(op *yadal.Operator) bool {
				return op.Metadata().Capability().Has(interfaces.Read, interfaces.Write)
			},
			tests: readWriteTests,
		},
	}
)

func getTests(name string) []testFunc {
	for _, test := range tests {
		if test.name == name {
			return test.tests
		}
	}
	log.Println("failed to find tests, please check the name")
	return nil
}

type builderFunc func() interfaces.Accessor

var (
	buildMap = map[string]builderFunc{
		// NOTES: uses UPPER word
		"FS": func() interfaces.Accessor {
			return fs.NewDriver(fs.Options{Root: os.Getenv("DAL_FS_ROOT")})
		},
		"S3": func() interfaces.Accessor {
			acc, err := s3.NewDriver(context.TODO(), s3.Options{
				Bucket:                     os.Getenv("DAL_S3_BUCKET"),
				Endpoint:                   os.Getenv("DAL_S3_ENDPOINT"),
				Root:                       os.Getenv("DAL_S3_ROOT"),
				Region:                     os.Getenv("DAL_S3_REGION"),
				AccessKey:                  os.Getenv("DAL_S3_ACCESS_KEY_ID"),
				SecretKey:                  os.Getenv("DAL_S3_SECRET_ACCESS_KEY"),
				SSEncryption:               nil,
				SSEncryptionAwsKmsKeyId:    nil,
				SSEncryptionCustomerAlgo:   nil,
				SSEncryptionCustomerKey:    nil,
				SSEncryptionCustomerKeyMD5: nil,
				EnableVirtualHostStyle:     false,
			})
			if err != nil {
				log.Fatal(err)
			}
			return acc
		},
	}
	s *zap.SugaredLogger
)

func init() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Println(err)
	}
	s = l.Sugar()
}

func setProviders() []*yadal.Operator {
	var op []*yadal.Operator
	debug := os.Getenv("TEST_DEBUG") == "on"
	log.Printf("Debug: %v\n", debug)

	var logging interfaces.Layer
	if debug {
		logging = layers.NewLoggingLayer(
			layers.SetLogger(
				layers.NewLoggerAdapter(s.Info, s.Infof),
			),
		)
	}

	for _, provider := range providers {
		PROVIDER := strings.ToUpper(provider)
		if os.Getenv(fmt.Sprintf("DAL_%s_TEST", PROVIDER)) == "on" {
			o := yadal.NewOperatorFromAccessor(buildMap[PROVIDER]())
			if debug {
				o.Layer(logging)
			}
			op = append(op, &o)
		}
	}
	return op
}

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %s", err)
	}
	os.Exit(m.Run())
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func runTests(t *testing.T, ops []*yadal.Operator, name string) {
	testList := getTests(name)
	log.Printf("---------------- running tests: %s ----------------", name)
	for _, operator := range ops {
		log.Printf("---------------- provider: %s ----------------", operator.Metadata().Provider().String())
		for _, test := range testList {
			log.Printf("---------------- running: %s ----------------", getFunctionName(test))
			test(t, operator)
		}
	}
}

func TestRun(t *testing.T) {
	p := setProviders()
	t.Run("basic", func(t *testing.T) {
		runTests(t, p, "basic")
	})
	t.Run("readWrite", func(t *testing.T) {
		runTests(t, p, "readWrite")
	})
}
