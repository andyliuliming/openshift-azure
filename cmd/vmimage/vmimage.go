package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/openshift/openshift-azure/pkg/tls"
	"github.com/openshift/openshift-azure/pkg/util/azureclient"
	"github.com/openshift/openshift-azure/pkg/util/log"
	"github.com/openshift/openshift-azure/pkg/util/random"
	"github.com/openshift/openshift-azure/pkg/util/resourceid"
	"github.com/openshift/openshift-azure/pkg/vmimage"
)

var (
	gitCommit = "unknown"

	timestamp = time.Now().UTC().Format("200601021504")

	logLevel                 = flag.String("loglevel", "Debug", "Valid values are Debug, Info, Warning, Error")
	location                 = flag.String("location", "eastus", "location")
	buildResourceGroup       = flag.String("buildResourceGroup", "vmimage-"+timestamp, "build resource group")
	deleteBuildResourceGroup = flag.Bool("deleteBuildResourceGroup", true, "delete build resource group after build")
	image                    = flag.String("image", "rhel7-3.11-"+timestamp, "image name")
	imageResourceGroup       = flag.String("imageResourceGroup", "images", "image resource group")
	imageStorageAccount      = flag.String("imageStorageAccount", "openshiftimages", "image storage account")
	imageContainer           = flag.String("imageContainer", "images", "image container")
	clientKey                = flag.String("clientKey", "secrets/client-key.pem", "cdn client key")
	clientCert               = flag.String("clientCert", "secrets/client-cert.pem", "cdn client cert")
)

func run(ctx context.Context, log *logrus.Entry) error {
	b, err := ioutil.ReadFile(*clientKey)
	if err != nil {
		return err
	}

	clientKey, err := tls.ParsePrivateKey(b)
	if err != nil {
		return err
	}

	b, err = ioutil.ReadFile(*clientCert)
	if err != nil {
		return err
	}

	clientCert, err := tls.ParseCert(b)
	if err != nil {
		return err
	}

	authorizer, err := azureclient.NewAuthorizerFromEnvironment("")
	if err != nil {
		return err
	}

	sshkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	b, err = tls.PrivateKeyAsBytes(sshkey)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("id_rsa", b, 0600)
	if err != nil {
		return err
	}

	domainNameLabel, err := random.LowerCaseAlphaString(20)
	if err != nil {
		return err
	}

	builder := vmimage.Builder{
		GitCommit:                gitCommit,
		Log:                      log,
		Deployments:              azureclient.NewDeploymentsClient(ctx, os.Getenv("AZURE_SUBSCRIPTION_ID"), authorizer),
		Groups:                   azureclient.NewGroupsClient(ctx, os.Getenv("AZURE_SUBSCRIPTION_ID"), authorizer),
		SubscriptionID:           os.Getenv("AZURE_SUBSCRIPTION_ID"),
		Location:                 *location,
		BuildResourceGroup:       *buildResourceGroup,
		DeleteBuildResourceGroup: *deleteBuildResourceGroup,
		DomainNameLabel:          domainNameLabel,
		Image:                    *image,
		ImageResourceGroup:       *imageResourceGroup,
		ImageStorageAccount:      *imageStorageAccount,
		ImageContainer:           *imageContainer,
		SSHKey:                   sshkey,
		ClientKey:                clientKey,
		ClientCert:               clientCert,
	}
	err = builder.Run(ctx)
	if err != nil {
		return err
	}

	return os.Remove("id_rsa")
}

func main() {
	flag.Parse()
	logrus.SetLevel(log.SanitizeLogLevel(*logLevel))
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	log := logrus.NewEntry(logrus.StandardLogger())
	log.Printf("vmimage starting, git commit %s", gitCommit)

	err := run(context.Background(), log)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("built image %s", resourceid.ResourceID(os.Getenv("AZURE_SUBSCRIPTION_ID"), *imageResourceGroup, "providers/Microsoft.Compute/images", *image))
}
