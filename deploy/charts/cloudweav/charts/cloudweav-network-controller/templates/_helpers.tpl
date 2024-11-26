{{/*
Expand the name of the chart.
*/}}
{{- define "cloudweav-network-controller.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cloudweav-network-controller.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cloudweav-network-controller.labels" -}}
helm.sh/chart: {{ include "cloudweav-network-controller.chart" . }}
{{ include "cloudweav-network-controller.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: network
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cloudweav-network-controller.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cloudweav-network-controller.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "cloudweav-network-controller-manager.labels" -}}
helm.sh/chart: {{ include "cloudweav-network-controller.chart" . }}
{{ include "cloudweav-network-controller-manager.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: network
{{- end }}

{{- define "cloudweav-network-controller-manager.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cloudweav-network-controller.name" . }}-manager
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "cloudweav-network-webhook.labels" -}}
helm.sh/chart: {{ include "cloudweav-network-controller.chart" . }}
{{ include "cloudweav-network-webhook.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: network-webhook
{{- end }}

{{- define "cloudweav-network-webhook.selectorLabels" -}}
app.kubernetes.io/name: cloudweav-network-webhook
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
