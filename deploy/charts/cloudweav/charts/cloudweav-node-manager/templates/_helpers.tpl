{{/*
Expand the name of the chart.
*/}}
{{- define "cloudweav-node-manager.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "cloudweav-node-manager-webhook.name" -}}
{{- default "cloudweav-node-manager-webhook" | trunc 63 }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cloudweav-node-manager.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cloudweav-node-manager.labels" -}}
helm.sh/chart: {{ include "cloudweav-node-manager.chart" . }}
{{ include "cloudweav-node-manager.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: node-manager
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cloudweav-node-manager.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cloudweav-node-manager.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "cloudweav-node-manager-webhook.labels" -}}
helm.sh/chart: {{ include "cloudweav-node-manager.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: node-manager
{{- end }}

{{- define "cloudweav-node-manager-webhook.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cloudweav-node-manager-webhook.name" . }}
{{- end }}