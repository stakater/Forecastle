{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "forecastle.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Allow the release namespace to be overridden
*/}}
{{- define "forecastle.namespace" -}}
{{- default .Release.Namespace .Values.namespace -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "forecastle.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "forecastle.labels.selector" -}}
app: {{ template "forecastle.name" . }}
group: {{ .Values.forecastle.labels.group }}
provider: {{ .Values.forecastle.labels.provider }}
{{- end -}}

{{- define "forecastle.labels.stakater" -}}
{{ template "forecastle.labels.selector" . }}
version: "{{ .Values.forecastle.labels.version }}"
{{- end -}}

{{- define "forecastle.labels.chart" -}}
chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
release: {{ .Release.Name | quote }}
heritage: {{ .Release.Service | quote }}
{{- end -}}

{{/*
Return the appropriate apiVersion for deployment.
*/}}
{{- define "deployment.apiVersion" -}}
{{- if semverCompare ">=1.9-0" .Capabilities.KubeVersion.GitVersion -}}
{{- print "apps/v1" -}}
{{- else -}}
{{- print "extensions/v1beta1" -}}
{{- end -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for podDisruptionBudget.
*/}}
{{- define "podDisruptionBudget.apiVersion" -}}
{{- if semverCompare ">=1.21-0" .Capabilities.KubeVersion.Version -}}
{{- print "policy/v1" -}}
{{- else -}}
{{- print "policy/v1beta1" -}}
{{- end -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for networkPolicy.
*/}}
{{- define "networkPolicy.apiVersion" -}}
{{- if semverCompare ">=1.4-0, <1.7-0" .Capabilities.KubeVersion.Version -}}
{{- print "extensions/v1beta1" -}}
{{- else if semverCompare "^1.7-0" .Capabilities.KubeVersion.Version -}}
{{- print "networking.k8s.io/v1" -}}
{{- end -}}
{{- end -}}

{{/*
Returns oauth redirection annotation according to the release name
*/}}
{{- define "serviceaccount.redirectreference" -}}
serviceaccounts.openshift.io/oauth-redirectreference.primary: '{"kind":"OAuthRedirectReference","apiVersion":"v1","reference":{"kind":"Route","name":"{{ template "forecastle.name" . }}"}}'
{{- end -}}

{{/*
Validate the RBAC scope value
*/}}
{{- define "forecastle.rbac.validate" -}}
{{- $scope := .Values.forecastle.rbac.scope -}}
{{- if not (or (eq $scope "cluster") (eq $scope "namespace")) -}}
{{- fail (printf "forecastle.rbac.scope must be either \"cluster\" or \"namespace\", got %q" $scope) -}}
{{- end -}}
{{- end -}}

{{/*
Shared RBAC rules for the resources Forecastle discovers
*/}}
{{- define "forecastle.rbac.rules" -}}
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list"]
{{- if .Values.forecastle.route.enabled }}
- apiGroups: ["route.openshift.io"]
  resources: ["routes"]
  verbs: ["get", "list"]
{{- end }}
- apiGroups: ["gateway.networking.k8s.io"]
  resources: ["httproutes"]
  verbs: ["get", "list"]
- apiGroups: ["traefik.containo.us"]
  resources: ["ingressroutes"]
  verbs: ["get", "list"]
- apiGroups: ["traefik.io"]
  resources: ["ingressroutes"]
  verbs: ["get", "list"]
- apiGroups: ["forecastle.stakater.com"]
  resources: ["forecastleapps"]
  verbs: ["get", "list"]
{{- end -}}
