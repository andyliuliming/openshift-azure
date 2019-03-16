// Package api defines the external API for the plugin.
package api

import (
	"context"

	adminapi "github.com/openshift/openshift-azure/pkg/api/admin/api"
)

// ContextKey is a type for context property bag payload keys
type ContextKey string

const (
	ContextKeyClientAuthorizer      ContextKey = "ClientAuthorizer"
	ContextKeyVaultClientAuthorizer ContextKey = "VaultClientAuthorizer"
	ContextAcceptLanguages          ContextKey = "AcceptLanguages"
)

type PluginStep string

const (
	PluginStepDeploy                                  PluginStep = "Deploy"
	PluginStepInitializeUpdateBlob                    PluginStep = "InitializeUpdateBlob"
	PluginStepResetUpdateBlob                         PluginStep = "ResetUpdateBlob"
	PluginStepCheckEtcdBlobExists                     PluginStep = "CheckEtcdBlobExists"
	PluginStepClientCreation                          PluginStep = "ClientCreation"
	PluginStepEnrichFromVault                         PluginStep = "EnrichFromVault"
	PluginStepScaleSetDelete                          PluginStep = "ScaleSetDelete"
	PluginStepWriteConfigBlob                         PluginStep = "WriteConfigBlob"
	PluginCreateOrUpdateConfigStorageAccount          PluginStep = "CreateOrUpdateConfigStorageAccount"
	PluginStepGenerateARM                             PluginStep = "GenerateARM"
	PluginStepWaitForWaitForOpenShiftAPI              PluginStep = "WaitForOpenShiftAPI"
	PluginStepWaitForNodes                            PluginStep = "WaitForNodes"
	PluginStepWaitForSyncPod                          PluginStep = "WaitForSyncPod"
	PluginStepWaitForConsoleHealth                    PluginStep = "WaitForConsoleHealth"
	PluginStepWaitForAdminConsoleHealth               PluginStep = "WaitForAdminConsoleHealth"
	PluginStepWaitForInfraDaemonSets                  PluginStep = "WaitForInfraDaemonSets"
	PluginStepWaitForInfraStatefulSets                PluginStep = "WaitForInfraStatefulSets"
	PluginStepWaitForInfraDeployments                 PluginStep = "WaitForInfraDeployments"
	PluginStepUpdateMasterAgentPoolHashMasterScaleSet PluginStep = "UpdateMasterAgentPoolHashMasterScaleSet"
	PluginStepUpdateMasterAgentPoolHashSyncPod        PluginStep = "UpdateMasterAgentPoolHashSyncPod"
	PluginStepUpdateMasterAgentPoolListVMs            PluginStep = "UpdateMasterAgentPoolListVMs"
	PluginStepUpdateMasterAgentPoolReadBlob           PluginStep = "UpdateMasterAgentPoolReadBlob"
	PluginStepUpdateMasterAgentPoolDrain              PluginStep = "UpdateMasterAgentPoolDrain"
	PluginStepUpdateMasterAgentPoolDeallocate         PluginStep = "UpdateMasterAgentPoolDeallocate"
	PluginStepUpdateMasterAgentPoolUpdateVMs          PluginStep = "UpdateMasterAgentPoolUpdateVMs"
	PluginStepUpdateMasterAgentPoolReimage            PluginStep = "UpdateMasterAgentPoolReimage"
	PluginStepUpdateMasterAgentPoolStart              PluginStep = "UpdateMasterAgentPoolStart"
	PluginStepUpdateMasterAgentPoolWaitForReady       PluginStep = "UpdateMasterAgentPoolWaitForReady"
	PluginStepUpdateMasterAgentPoolUpdateBlob         PluginStep = "UpdateMasterAgentPoolUpdateBlob"
	PluginStepUpdateWorkerAgentPoolHashWorkerScaleSet PluginStep = "UpdateWorkerAgentPoolHashWorkerScaleSet"
	PluginStepUpdateWorkerAgentPoolListVMs            PluginStep = "UpdateWorkerAgentPoolListVMs"
	PluginStepUpdateWorkerAgentPoolListScaleSets      PluginStep = "UpdateWorkerAgentPoolListScaleSets"
	PluginStepUpdateWorkerAgentPoolReadBlob           PluginStep = "UpdateWorkerAgentPoolReadBlob"
	PluginStepUpdateWorkerAgentPoolDrain              PluginStep = "UpdateWorkerAgentPoolDrain"
	PluginStepUpdateWorkerAgentPoolCreateScaleSet     PluginStep = "UpdateWorkerAgentPoolCreateScaleSet"
	PluginStepUpdateWorkerAgentPoolUpdateScaleSet     PluginStep = "UpdateWorkerAgentPoolUpdateScaleSet"
	PluginStepUpdateWorkerAgentPoolDeleteScaleSet     PluginStep = "UpdateWorkerAgentPoolDeleteScaleSet"
	PluginStepUpdateWorkerAgentPoolWaitForReady       PluginStep = "UpdateWorkerAgentPoolWaitForReady"
	PluginStepUpdateWorkerAgentPoolUpdateBlob         PluginStep = "UpdateWorkerAgentPoolUpdateBlob"
	PluginStepUpdateWorkerAgentPoolDeleteVM           PluginStep = "UpdateWorkerAgentPoolDeleteVM"
	PluginStepUpdateSyncPodReadBlob                   PluginStep = "UpdateSyncPodReadBlob"
	PluginStepUpdateSyncPodHashSyncPod                PluginStep = "UpdateSyncPodHashSyncPod"
	PluginStepUpdateSyncPodDeletePod                  PluginStep = "UpdateSyncPodDeletePod"
	PluginStepUpdateSyncPodWaitForReady               PluginStep = "UpdateSyncPodWaitForReady"
	PluginStepUpdateSyncPodUpdateBlob                 PluginStep = "UpdateSyncPodUpdateBlob"
	PluginStepInvalidateClusterSecrets                PluginStep = "InvalidateClusterSecrets"
	PluginStepRegenerateClusterSecrets                PluginStep = "RegenerateClusterSecrets"
)

type Command string

const (
	CommandRestartNetworkManager = "RestartNetworkManager"
	CommandRestartKubelet        = "RestartKubelet"
)

// PluginError error returned by CreateOrUpdate to specify the step that failed.
type PluginError struct {
	Err  error
	Step PluginStep
}

var _ error = &PluginError{}

func (pe *PluginError) Error() string {
	return string(pe.Step) + ": " + pe.Err.Error()
}

// DeployFn makes it possible to plug in different logic to the deploy.
// The implementor must initiate a deployment of the given template using
// mode resources.Incremental and wait for it to complete.
type DeployFn func(context.Context, map[string]interface{}) error

// TestConfig holds all testing variables. It should be the zero value in
// production.
type TestConfig struct {
	RunningUnderTest   bool
	ImageResourceGroup string
	ImageResourceName  string
}

// Plugin is the main interface to openshift-azure
type Plugin interface {
	// Validate exists (a) to be able to place validation logic in a
	// single place in the event of multiple external API versions, and (b) to
	// be able to compare a new API manifest against a pre-existing API manifest
	// (for update, upgrade, etc.)
	// externalOnly indicates that fields set by the RP (FQDN and routerProfile.FQDN)
	// should be excluded.
	Validate(ctx context.Context, new, old *OpenShiftManagedCluster, externalOnly bool) []error

	// ValidateAdmin is used for validating admin API requests.
	ValidateAdmin(ctx context.Context, new, old *OpenShiftManagedCluster) []error

	// ValidatePluginTemplate validates external config request
	ValidatePluginTemplate(ctx context.Context) []error

	// GenerateConfig ensures all the necessary in-cluster config is generated
	// for an Openshift cluster.
	GenerateConfig(ctx context.Context, cs *OpenShiftManagedCluster) error

	// CreateOrUpdate either deploys or runs the update depending on the isUpdate argument
	// this will call the deployer.
	CreateOrUpdate(ctx context.Context, cs *OpenShiftManagedCluster, isUpdate bool, deployer DeployFn) *PluginError

	GenevaActions
}

// GenevaActions is the interface for all geneva actions
type GenevaActions interface {
	// RecoverEtcdCluster recovers the cluster's etcd using the backup specified in the pluginConfig
	RecoverEtcdCluster(ctx context.Context, cs *OpenShiftManagedCluster, deployer DeployFn, backupBlob string) *PluginError

	// RotateClusterSecrets rotates the secrets in a cluster's config blob and then updates the cluster
	RotateClusterSecrets(ctx context.Context, cs *OpenShiftManagedCluster, deployer DeployFn) *PluginError

	// GetControlPlanePods fetches a consolidated list of the control plane pods in the cluster
	GetControlPlanePods(ctx context.Context, oc *OpenShiftManagedCluster) ([]byte, error)

	// ForceUpdate forces rotates all vms in a cluster
	ForceUpdate(ctx context.Context, cs *OpenShiftManagedCluster, deployer DeployFn) *PluginError

	// ListClusterVMs returns the hostnames of all vms in a cluster
	ListClusterVMs(ctx context.Context, cs *OpenShiftManagedCluster) (*adminapi.GenevaActionListClusterVMs, error)

	// Reimage reimages a virtual machine in the cluster
	Reimage(ctx context.Context, oc *OpenShiftManagedCluster, hostname string) error

	// BackupEtcdCluster backs up the cluster's etcd
	BackupEtcdCluster(ctx context.Context, cs *OpenShiftManagedCluster, backupName string) error

	// RunCommand runs a predefined command on a virtual machine in the cluster
	RunCommand(ctx context.Context, cs *OpenShiftManagedCluster, hostname string, command Command) error

	// GetPluginVersion fetches the RP plugin version
	GetPluginVersion(ctx context.Context) *adminapi.GenevaActionPluginVersion
}
