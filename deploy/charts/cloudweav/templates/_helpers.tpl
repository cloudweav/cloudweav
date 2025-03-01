{{/* vim: set filetype=mustache: */}}

{{/*
Expand the name of the chart, which is no longer than 63 chars.
We truncate at 63 chars as some Kubernetes naming policies are limited to this.
*/}}
{{- define "cloudweav.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{/*
Create a default fully qualified app name, which is no longer than 63 chars.
We truncate at 63 chars as some Kubernetes naming policies are limited to this.
*/}}
{{- define "cloudweav.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cloudweav.chartref" -}}
{{- replace "+" "_" .Chart.Version | printf "%s-%s" .Chart.Name -}}
{{- end }}

{{/*
Generate immutable labels.
We use the immutable labels to select low-level components(like the Deployment selects the Pods).
*/}}
{{- define "cloudweav.immutableLabels" -}}
helm.sh/release: {{ .Release.Name }}
app.kubernetes.io/part-of: {{ template "cloudweav.name" . }}
{{- end }}

{{/*
Generate basic labels.
*/}}
{{- define "cloudweav.labels" -}}
app.kubernetes.io/managed-by: {{ default "helm" $.Release.Service | quote }}
helm.sh/chart: {{ template "cloudweav.chartref" . }}
app.kubernetes.io/version: {{ $.Chart.AppVersion | quote }}
{{ include "cloudweav.immutableLabels" . }}
{{- end }}

{{/*
Generate API affinity. It makes pods of a workoad to run on different nodes.
*/}}
{{- define "cloudweav.apiAffinity" -}}
podAntiAffinity:
  requiredDuringSchedulingIgnoredDuringExecution:
    - labelSelector:
        matchExpressions:
          - key: app.kubernetes.io/name
            operator: In
            values:
              - cloudweav
          - key: app.kubernetes.io/component
            operator: In
            values:
              - {{ .component }}
          - key: app.kubernetes.io/version
            operator: In
            values:
              - {{ .root.Chart.AppVersion }}
      topologyKey: kubernetes.io/hostname
nodeAffinity:
  requiredDuringSchedulingIgnoredDuringExecution:
    nodeSelectorTerms:
      - matchExpressions:
          - key: kubernetes.io/os
            operator: In
            values:
              - linux
{{- end }}

{{/*
NB(thxCode): Use this value to unify the control tag and condition of KubeVirt Operator.
*/}}
{{- define "conditions.is_kubevirt_operator_enabled" }}
{{- $kubevirtOperatorEnabled := (index .Values "kubevirt-operator" "enabled") | toString -}}
{{- if ne $kubevirtOperatorEnabled "<nil>" -}}
{{- $kubevirtOperatorEnabled -}}
{{- else -}}
{{- .Values.tags.kubevirt | toString -}}
{{- end -}}
{{- end }}

{{/*
NB(thxCode): Use this value to unify the control tag and condition of KubeVirt.
*/}}
{{- define "conditions.is_kubevirt_enabled" }}
{{- $kubevirtEnabled := (index .Values "kubevirt" "enabled") | toString -}}
{{- if ne $kubevirtEnabled "<nil>" }}
{{- $kubevirtEnabled -}}
{{- else -}}
{{- .Values.tags.kubevirt | toString -}}
{{- end -}}
{{- end }}

{{/*
Get Support-bundle-kit image environment for updating the default values per current release.
*/}}
{{- define "cloudweav.supportBundleImageEnv" -}}
{{- $result := dict -}}
{{- range $k, $v := .Values -}}
{{- if eq (toString $k) "support-bundle-kit" -}}
{{- $result = $v -}}
{{- end -}}
{{- end -}}
{{- with $result -}}
{{- with .image -}}
- name: CLOUDWEAV_SUPPORT_BUNDLE_IMAGE_DEFAULT_VALUE
  value: {{ printf "{\"repository\":\"%s\",\"tag\":\"%s\",\"imagePullPolicy\":\"%s\"}" .repository .tag .imagePullPolicy | squote }}
{{- end -}}
{{- end -}}
{{- end }}
