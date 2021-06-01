package main

import (
	"bytes"
	"context"
	"encoding/json"

	//"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerregistry/mgmt/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	ShellScriptFileName string = "getimagesha.sh"
)

var (
	debug  = pflag.Bool("debug", true, "sets log to debug level")
	server *Server
	ctx    context.Context
)

// LogHook is used to setup custom hooks
type LogHook struct {
	Writer    io.Writer
	Loglevels []log.Level
}

// Image struct/
type Image struct {
	registry string
	repo     string
	tag      string
	digest   string
}

func main() {
	//pflag.Parse()

	var err error

	setupLogger()

	ctx = context.Background()
	server, err = NewServer()
	if err != nil {
		log.Fatalf("[error] : %v", err)
	}
	http.HandleFunc("/process", handle)
	http.ListenAndServe(":8090", nil)

	os.Exit(0)
}

func handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	imageStr := req.URL.Query().Get("image") // e.g. : oss/kubernetes/aks/etcd-operator
	if imageStr == "" {
		log.Info("Failed to provide image to query")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(nil)
		return
	}
	image := newImage(imageStr)
	if strings.EqualFold(image.digest, "") {
		log.Errorf("Digest is empty.")
		return
	}

	scanInfo, err := server.Process(ctx, image)
	if err != nil {
		log.Infof("[error] : %s", err)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(scanInfo)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(scanInfo)
	}

}

func newImage(imageStr string) (image Image) {
	// e.g image fields:
	// registry := "upstream.azurecr.io"
	// repo := "oss/kubernetes/ingress/nginx-ingress-controller"
	// tag := "0.16.2"
	registry := strings.Split(imageStr, "/")[0]
	repo := strings.Replace(imageStr, registry+"/", "", 1)
	tag := "latest"
	if strings.Contains(repo, ":") {
		tag = strings.Split(repo, ":")[1]
		repo = strings.Replace(repo, ":"+tag, "", 1)
	}
	image = Image{registry: registry, repo: repo, tag: tag, digest: ""}
	// In case that the image was deployed with the digest:
	if strings.Contains(imageStr, "@sha256") {
		image.digest = strings.Split(imageStr, "@")[1]
		// else, extract digest first.
	} else {
		image.digest = tag2Digest(image)
	}
	return image
}

// Convert tag to digest - using shell script (getimagesha.sh)
func tag2Digest(image Image) (digest string) {
	if image.registry == "" || image.repo == "" || image.tag == "" {
		log.Errorf("Invalid image - registry or repo or tag are empty : registry = %s, repo = %s, tag = %s", image.registry, image.repo, image.tag)
		return ""
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(
		"sh",
		ShellScriptFileName,
		image.registry,
		image.repo,
		image.tag,
	)

	log.Infof("cmd: %v", cmd)
	cmd.Dir = dir
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stderr, cmd.Stdout = stderr, stdout

	err = cmd.Run()
	output := stdout.String()
	log.Infof("output: %s", output)
	if err != nil {
		log.Errorf("error invoking cmd, err: %v, output: %v", err, stderr.String())
	}

	if output == "null\n" {
		log.Infof("[error] : could not find valid digest %s", output)
	} else {
		digest = strings.TrimSuffix(output, "\n")
	}
	//TODO Check that the digest is valid.
	image.digest = digest
	log.Infof("Digest successfully extracted: %s", digest)
	return digest
}

// setupLogger sets up hooks to redirect stdout and stderr
func setupLogger() {
	log.SetOutput(ioutil.Discard)

	// set log level
	log.SetLevel(log.InfoLevel)
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	// add hook to send info, debug, warn level logs to stdout
	log.AddHook(&LogHook{
		Writer: os.Stdout,
		Loglevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
			log.WarnLevel,
		},
	})

	// add hook to send panic, fatal, error logs to stderr
	log.AddHook(&LogHook{
		Writer: os.Stderr,
		Loglevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		},
	})
}

// Fire is called when logging function with current hook is called
// write to appropriate writer based on log level
func (hook *LogHook) Fire(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels defines log levels at which hook is triggered
func (hook *LogHook) Levels() []log.Level {
	return hook.Loglevels
}

func tag2DigestSDK(image Image) {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	containerRegistryClient := containerregistry.NewRegistriesClient(os.Getenv("SubscriptionID"))
	containerRegistryClient.Authorizer = authorizer
	containerRegistryCredential, err := containerRegistryClient.ListCredentials(context.Background(), os.Getenv("ResourceGroupName"), os.Getenv("Registry"))

	basicAuthorizer := autorest.NewBasicAuthorizer("<userName>", "<password>")
	tagClient := repository.NewTagClient("https://" + "<loginServer>")
	tagClient.Authorizer = basicAuthorizer
	repositoryTagDetails, err := tagClient.GetList(context.Background(), " <repositoryName>", "", nil, "", "")
}