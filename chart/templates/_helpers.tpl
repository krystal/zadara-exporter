{{- define "zadaraexporter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "zadaraexporter.determineFullname" -}}
{{- if contains .ChartName .ReleaseName }}
{{- .ReleaseName | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .ReleaseName .ChartName | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{- define "zadaraexporter.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- include "zadaraexporter.determineFullname" (dict "ChartName" $name "ReleaseName" .Release.Name) }}
{{- end }}
{{- end }}

{{- define "zadaraexporter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "zadaraexporter.labels" -}}
helm.sh/chart: {{ include "zadaraexporter.chart" . }}
{{ include "zadaraexporter.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "zadaraexporter.selectorLabels" -}}
app.kubernetes.io/name: {{ include "zadaraexporter.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
