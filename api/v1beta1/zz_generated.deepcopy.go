// +build !ignore_autogenerated

/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Admin) DeepCopyInto(out *Admin) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Admin.
func (in *Admin) DeepCopy() *Admin {
	if in == nil {
		return nil
	}
	out := new(Admin)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ambassador) DeepCopyInto(out *Ambassador) {
	*out = *in
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(AmbassadorService)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ambassador.
func (in *Ambassador) DeepCopy() *Ambassador {
	if in == nil {
		return nil
	}
	out := new(Ambassador)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AmbassadorService) DeepCopyInto(out *AmbassadorService) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]*Port, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Port)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AmbassadorService.
func (in *AmbassadorService) DeepCopy() *AmbassadorService {
	if in == nil {
		return nil
	}
	out := new(AmbassadorService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Argocd) DeepCopyInto(out *Argocd) {
	*out = *in
	if in.CustomImage != nil {
		in, out := &in.CustomImage, &out.CustomImage
		*out = new(ArgocdCustomImage)
		(*in).DeepCopyInto(*out)
	}
	if in.Network != nil {
		in, out := &in.Network, &out.Network
		*out = new(Network)
		**out = **in
	}
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(ArgocdAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.Rbac != nil {
		in, out := &in.Rbac, &out.Rbac
		*out = new(Rbac)
		**out = **in
	}
	if in.Repositories != nil {
		in, out := &in.Repositories, &out.Repositories
		*out = make([]*ArgocdRepository, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(ArgocdRepository)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Argocd.
func (in *Argocd) DeepCopy() *Argocd {
	if in == nil {
		return nil
	}
	out := new(Argocd)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdAuth) DeepCopyInto(out *ArgocdAuth) {
	*out = *in
	if in.OIDC != nil {
		in, out := &in.OIDC, &out.OIDC
		*out = new(ArgocdOIDC)
		(*in).DeepCopyInto(*out)
	}
	if in.GithubConnector != nil {
		in, out := &in.GithubConnector, &out.GithubConnector
		*out = new(ArgocdGithubConnector)
		(*in).DeepCopyInto(*out)
	}
	if in.GitlabConnector != nil {
		in, out := &in.GitlabConnector, &out.GitlabConnector
		*out = new(ArgocdGitlabConnector)
		(*in).DeepCopyInto(*out)
	}
	if in.GoogleConnector != nil {
		in, out := &in.GoogleConnector, &out.GoogleConnector
		*out = new(ArgocdGoogleConnector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdAuth.
func (in *ArgocdAuth) DeepCopy() *ArgocdAuth {
	if in == nil {
		return nil
	}
	out := new(ArgocdAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdClaim) DeepCopyInto(out *ArgocdClaim) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdClaim.
func (in *ArgocdClaim) DeepCopy() *ArgocdClaim {
	if in == nil {
		return nil
	}
	out := new(ArgocdClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdCustomImage) DeepCopyInto(out *ArgocdCustomImage) {
	*out = *in
	if in.GopassSSHKey != nil {
		in, out := &in.GopassSSHKey, &out.GopassSSHKey
		*out = new(ArgocdSecret)
		**out = **in
	}
	if in.GopassGPGKey != nil {
		in, out := &in.GopassGPGKey, &out.GopassGPGKey
		*out = new(ArgocdSecret)
		**out = **in
	}
	if in.GopassStores != nil {
		in, out := &in.GopassStores, &out.GopassStores
		*out = make([]*ArgocdGopassStore, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(ArgocdGopassStore)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdCustomImage.
func (in *ArgocdCustomImage) DeepCopy() *ArgocdCustomImage {
	if in == nil {
		return nil
	}
	out := new(ArgocdCustomImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdGithubConfig) DeepCopyInto(out *ArgocdGithubConfig) {
	*out = *in
	if in.Orgs != nil {
		in, out := &in.Orgs, &out.Orgs
		*out = make([]*ArgocdGithubOrg, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(ArgocdGithubOrg)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdGithubConfig.
func (in *ArgocdGithubConfig) DeepCopy() *ArgocdGithubConfig {
	if in == nil {
		return nil
	}
	out := new(ArgocdGithubConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdGithubConnector) DeepCopyInto(out *ArgocdGithubConnector) {
	*out = *in
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(ArgocdGithubConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdGithubConnector.
func (in *ArgocdGithubConnector) DeepCopy() *ArgocdGithubConnector {
	if in == nil {
		return nil
	}
	out := new(ArgocdGithubConnector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdGithubOrg) DeepCopyInto(out *ArgocdGithubOrg) {
	*out = *in
	if in.Teams != nil {
		in, out := &in.Teams, &out.Teams
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdGithubOrg.
func (in *ArgocdGithubOrg) DeepCopy() *ArgocdGithubOrg {
	if in == nil {
		return nil
	}
	out := new(ArgocdGithubOrg)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdGitlabConfig) DeepCopyInto(out *ArgocdGitlabConfig) {
	*out = *in
	if in.Groups != nil {
		in, out := &in.Groups, &out.Groups
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdGitlabConfig.
func (in *ArgocdGitlabConfig) DeepCopy() *ArgocdGitlabConfig {
	if in == nil {
		return nil
	}
	out := new(ArgocdGitlabConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdGitlabConnector) DeepCopyInto(out *ArgocdGitlabConnector) {
	*out = *in
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(ArgocdGitlabConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdGitlabConnector.
func (in *ArgocdGitlabConnector) DeepCopy() *ArgocdGitlabConnector {
	if in == nil {
		return nil
	}
	out := new(ArgocdGitlabConnector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdGoogleConfig) DeepCopyInto(out *ArgocdGoogleConfig) {
	*out = *in
	if in.HostedDomains != nil {
		in, out := &in.HostedDomains, &out.HostedDomains
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Groups != nil {
		in, out := &in.Groups, &out.Groups
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdGoogleConfig.
func (in *ArgocdGoogleConfig) DeepCopy() *ArgocdGoogleConfig {
	if in == nil {
		return nil
	}
	out := new(ArgocdGoogleConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdGoogleConnector) DeepCopyInto(out *ArgocdGoogleConnector) {
	*out = *in
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(ArgocdGoogleConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdGoogleConnector.
func (in *ArgocdGoogleConnector) DeepCopy() *ArgocdGoogleConnector {
	if in == nil {
		return nil
	}
	out := new(ArgocdGoogleConnector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdGopassStore) DeepCopyInto(out *ArgocdGopassStore) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdGopassStore.
func (in *ArgocdGopassStore) DeepCopy() *ArgocdGopassStore {
	if in == nil {
		return nil
	}
	out := new(ArgocdGopassStore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdOIDC) DeepCopyInto(out *ArgocdOIDC) {
	*out = *in
	if in.RequestedScopes != nil {
		in, out := &in.RequestedScopes, &out.RequestedScopes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RequestedIDTokenClaims != nil {
		in, out := &in.RequestedIDTokenClaims, &out.RequestedIDTokenClaims
		*out = make(map[string]ArgocdClaim, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdOIDC.
func (in *ArgocdOIDC) DeepCopy() *ArgocdOIDC {
	if in == nil {
		return nil
	}
	out := new(ArgocdOIDC)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdRepoSecret) DeepCopyInto(out *ArgocdRepoSecret) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdRepoSecret.
func (in *ArgocdRepoSecret) DeepCopy() *ArgocdRepoSecret {
	if in == nil {
		return nil
	}
	out := new(ArgocdRepoSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdRepository) DeepCopyInto(out *ArgocdRepository) {
	*out = *in
	if in.UsernameSecret != nil {
		in, out := &in.UsernameSecret, &out.UsernameSecret
		*out = new(ArgocdRepoSecret)
		**out = **in
	}
	if in.PasswordSecret != nil {
		in, out := &in.PasswordSecret, &out.PasswordSecret
		*out = new(ArgocdRepoSecret)
		**out = **in
	}
	if in.CertificateSecret != nil {
		in, out := &in.CertificateSecret, &out.CertificateSecret
		*out = new(ArgocdRepoSecret)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdRepository.
func (in *ArgocdRepository) DeepCopy() *ArgocdRepository {
	if in == nil {
		return nil
	}
	out := new(ArgocdRepository)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArgocdSecret) DeepCopyInto(out *ArgocdSecret) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArgocdSecret.
func (in *ArgocdSecret) DeepCopy() *ArgocdSecret {
	if in == nil {
		return nil
	}
	out := new(ArgocdSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Datasource) DeepCopyInto(out *Datasource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Datasource.
func (in *Datasource) DeepCopy() *Datasource {
	if in == nil {
		return nil
	}
	out := new(Datasource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Grafana) DeepCopyInto(out *Grafana) {
	*out = *in
	if in.Admin != nil {
		in, out := &in.Admin, &out.Admin
		*out = new(Admin)
		**out = **in
	}
	if in.Datasources != nil {
		in, out := &in.Datasources, &out.Datasources
		*out = make([]*Datasource, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Datasource)
				**out = **in
			}
		}
	}
	if in.DashboardProviders != nil {
		in, out := &in.DashboardProviders, &out.DashboardProviders
		*out = make([]*Provider, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Provider)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Storage != nil {
		in, out := &in.Storage, &out.Storage
		*out = new(StorageSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Network != nil {
		in, out := &in.Network, &out.Network
		*out = new(Network)
		**out = **in
	}
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(GrafanaAuth)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Grafana.
func (in *Grafana) DeepCopy() *Grafana {
	if in == nil {
		return nil
	}
	out := new(Grafana)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrafanaAuth) DeepCopyInto(out *GrafanaAuth) {
	*out = *in
	if in.Google != nil {
		in, out := &in.Google, &out.Google
		*out = new(GrafanaGoogleAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.Github != nil {
		in, out := &in.Github, &out.Github
		*out = new(GrafanaGithubAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.Gitlab != nil {
		in, out := &in.Gitlab, &out.Gitlab
		*out = new(GrafanaGitlabAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.GenericOAuth != nil {
		in, out := &in.GenericOAuth, &out.GenericOAuth
		*out = new(GrafanaGenericOAuth)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrafanaAuth.
func (in *GrafanaAuth) DeepCopy() *GrafanaAuth {
	if in == nil {
		return nil
	}
	out := new(GrafanaAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrafanaGenericOAuth) DeepCopyInto(out *GrafanaGenericOAuth) {
	*out = *in
	if in.Scopes != nil {
		in, out := &in.Scopes, &out.Scopes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AllowedDomains != nil {
		in, out := &in.AllowedDomains, &out.AllowedDomains
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrafanaGenericOAuth.
func (in *GrafanaGenericOAuth) DeepCopy() *GrafanaGenericOAuth {
	if in == nil {
		return nil
	}
	out := new(GrafanaGenericOAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrafanaGithubAuth) DeepCopyInto(out *GrafanaGithubAuth) {
	*out = *in
	if in.AllowedOrganizations != nil {
		in, out := &in.AllowedOrganizations, &out.AllowedOrganizations
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.TeamIDs != nil {
		in, out := &in.TeamIDs, &out.TeamIDs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrafanaGithubAuth.
func (in *GrafanaGithubAuth) DeepCopy() *GrafanaGithubAuth {
	if in == nil {
		return nil
	}
	out := new(GrafanaGithubAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrafanaGitlabAuth) DeepCopyInto(out *GrafanaGitlabAuth) {
	*out = *in
	if in.AllowedGroups != nil {
		in, out := &in.AllowedGroups, &out.AllowedGroups
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrafanaGitlabAuth.
func (in *GrafanaGitlabAuth) DeepCopy() *GrafanaGitlabAuth {
	if in == nil {
		return nil
	}
	out := new(GrafanaGitlabAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrafanaGoogleAuth) DeepCopyInto(out *GrafanaGoogleAuth) {
	*out = *in
	if in.AllowedDomains != nil {
		in, out := &in.AllowedDomains, &out.AllowedDomains
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrafanaGoogleAuth.
func (in *GrafanaGoogleAuth) DeepCopy() *GrafanaGoogleAuth {
	if in == nil {
		return nil
	}
	out := new(GrafanaGoogleAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeStateMetrics) DeepCopyInto(out *KubeStateMetrics) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeStateMetrics.
func (in *KubeStateMetrics) DeepCopy() *KubeStateMetrics {
	if in == nil {
		return nil
	}
	out := new(KubeStateMetrics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoggingOperator) DeepCopyInto(out *LoggingOperator) {
	*out = *in
	if in.FluentdPVC != nil {
		in, out := &in.FluentdPVC, &out.FluentdPVC
		*out = new(StorageSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoggingOperator.
func (in *LoggingOperator) DeepCopy() *LoggingOperator {
	if in == nil {
		return nil
	}
	out := new(LoggingOperator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Logs) DeepCopyInto(out *Logs) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Logs.
func (in *Logs) DeepCopy() *Logs {
	if in == nil {
		return nil
	}
	out := new(Logs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Loki) DeepCopyInto(out *Loki) {
	*out = *in
	if in.Logs != nil {
		in, out := &in.Logs, &out.Logs
		*out = new(Logs)
		**out = **in
	}
	if in.Storage != nil {
		in, out := &in.Storage, &out.Storage
		*out = new(StorageSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Loki.
func (in *Loki) DeepCopy() *Loki {
	if in == nil {
		return nil
	}
	out := new(Loki)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Metrics) DeepCopyInto(out *Metrics) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Metrics.
func (in *Metrics) DeepCopy() *Metrics {
	if in == nil {
		return nil
	}
	out := new(Metrics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Network) DeepCopyInto(out *Network) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Network.
func (in *Network) DeepCopy() *Network {
	if in == nil {
		return nil
	}
	out := new(Network)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Port) DeepCopyInto(out *Port) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Port.
func (in *Port) DeepCopy() *Port {
	if in == nil {
		return nil
	}
	out := new(Port)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostApply) DeepCopyInto(out *PostApply) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostApply.
func (in *PostApply) DeepCopy() *PostApply {
	if in == nil {
		return nil
	}
	out := new(PostApply)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreApply) DeepCopyInto(out *PreApply) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreApply.
func (in *PreApply) DeepCopy() *PreApply {
	if in == nil {
		return nil
	}
	out := new(PreApply)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Prometheus) DeepCopyInto(out *Prometheus) {
	*out = *in
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = new(Metrics)
		**out = **in
	}
	if in.Storage != nil {
		in, out := &in.Storage, &out.Storage
		*out = new(StorageSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Prometheus.
func (in *Prometheus) DeepCopy() *Prometheus {
	if in == nil {
		return nil
	}
	out := new(Prometheus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusNodeExporter) DeepCopyInto(out *PrometheusNodeExporter) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusNodeExporter.
func (in *PrometheusNodeExporter) DeepCopy() *PrometheusNodeExporter {
	if in == nil {
		return nil
	}
	out := new(PrometheusNodeExporter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusOperator) DeepCopyInto(out *PrometheusOperator) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusOperator.
func (in *PrometheusOperator) DeepCopy() *PrometheusOperator {
	if in == nil {
		return nil
	}
	out := new(PrometheusOperator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Provider) DeepCopyInto(out *Provider) {
	*out = *in
	if in.ConfigMaps != nil {
		in, out := &in.ConfigMaps, &out.ConfigMaps
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Provider.
func (in *Provider) DeepCopy() *Provider {
	if in == nil {
		return nil
	}
	out := new(Provider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rbac) DeepCopyInto(out *Rbac) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rbac.
func (in *Rbac) DeepCopy() *Rbac {
	if in == nil {
		return nil
	}
	out := new(Rbac)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageSpec) DeepCopyInto(out *StorageSpec) {
	*out = *in
	if in.AccessModes != nil {
		in, out := &in.AccessModes, &out.AccessModes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageSpec.
func (in *StorageSpec) DeepCopy() *StorageSpec {
	if in == nil {
		return nil
	}
	out := new(StorageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Toolset) DeepCopyInto(out *Toolset) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(ToolsetSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(ToolsetStatus)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Toolset.
func (in *Toolset) DeepCopy() *Toolset {
	if in == nil {
		return nil
	}
	out := new(Toolset)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Toolset) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ToolsetList) DeepCopyInto(out *ToolsetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]*Toolset, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Toolset)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ToolsetList.
func (in *ToolsetList) DeepCopy() *ToolsetList {
	if in == nil {
		return nil
	}
	out := new(ToolsetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ToolsetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ToolsetSpec) DeepCopyInto(out *ToolsetSpec) {
	*out = *in
	if in.PreApply != nil {
		in, out := &in.PreApply, &out.PreApply
		*out = new(PreApply)
		**out = **in
	}
	if in.PostApply != nil {
		in, out := &in.PostApply, &out.PostApply
		*out = new(PostApply)
		**out = **in
	}
	if in.PrometheusOperator != nil {
		in, out := &in.PrometheusOperator, &out.PrometheusOperator
		*out = new(PrometheusOperator)
		**out = **in
	}
	if in.LoggingOperator != nil {
		in, out := &in.LoggingOperator, &out.LoggingOperator
		*out = new(LoggingOperator)
		(*in).DeepCopyInto(*out)
	}
	if in.PrometheusNodeExporter != nil {
		in, out := &in.PrometheusNodeExporter, &out.PrometheusNodeExporter
		*out = new(PrometheusNodeExporter)
		**out = **in
	}
	if in.Grafana != nil {
		in, out := &in.Grafana, &out.Grafana
		*out = new(Grafana)
		(*in).DeepCopyInto(*out)
	}
	if in.Ambassador != nil {
		in, out := &in.Ambassador, &out.Ambassador
		*out = new(Ambassador)
		(*in).DeepCopyInto(*out)
	}
	if in.KubeStateMetrics != nil {
		in, out := &in.KubeStateMetrics, &out.KubeStateMetrics
		*out = new(KubeStateMetrics)
		**out = **in
	}
	if in.Argocd != nil {
		in, out := &in.Argocd, &out.Argocd
		*out = new(Argocd)
		(*in).DeepCopyInto(*out)
	}
	if in.Prometheus != nil {
		in, out := &in.Prometheus, &out.Prometheus
		*out = new(Prometheus)
		(*in).DeepCopyInto(*out)
	}
	if in.Loki != nil {
		in, out := &in.Loki, &out.Loki
		*out = new(Loki)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ToolsetSpec.
func (in *ToolsetSpec) DeepCopy() *ToolsetSpec {
	if in == nil {
		return nil
	}
	out := new(ToolsetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ToolsetStatus) DeepCopyInto(out *ToolsetStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ToolsetStatus.
func (in *ToolsetStatus) DeepCopy() *ToolsetStatus {
	if in == nil {
		return nil
	}
	out := new(ToolsetStatus)
	in.DeepCopyInto(out)
	return out
}
