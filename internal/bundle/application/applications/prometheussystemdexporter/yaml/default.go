package yaml

func GetDefault() interface{} {
	return map[string]interface{}{
		"kind":       "Deployment",
		"apiVersion": "apps/v1",
		"metadata": map[string]interface{}{
			"name":      "systemd-exporter",
			"namespace": "caos-system",
			"labels": map[string]string{
				"app": "systemd-exporter",
			},
		},
		"spec": map[string]interface{}{
			"selector": map[string]interface{}{
				"matchLabels": map[string]string{
					"app": "systemd-exporter",
				},
			},
			"updateStrategy": map[string]interface{}{
				"rollingUpdate": map[string]string{
					"maxUnavailable": "100%",
				},
				"type": "RollingUpdate",
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]string{
						"app": "systemd-exporter",
					},
					"annotations": map[string]string{
						"prometheus.io/scrape": "true",
						"prometheus.io/path":   "/metrics",
						"prometheus.io/port":   "9558",
					},
				},
				"spec": map[string]interface{}{
					"tolerations": []map[string]string{
						{"key": "node-role.kubernetes.io/master"},
						{"operator": "Equal"},
						{"effect": "NoSchedule"},
					},
					"securityContext": map[string]uint8{
						"runAsUser": 0,
					},
					"hostPID": true,
					"containers": []map[string]interface{}{
						{"name": "systemd-exporter"},
						{"image": "quay.io/povilasv/systemd_exporter:v0.2.0"},
						{"securityContext": map[string]bool{
							"privileged": true,
						}},
						{"args": []string{
							"--log.level=info",
							"--path.procfs=/host/proc",
							"--web.disable-exporter-metrics",
							"--collector.unit-whitelist=kubelet.service|docker.service|node-agentd.service|firewalld.service|keepalived.service|nginx.service|sshd.service",
						}},
						{"ports": []map[string]interface{}{
							{"name": "metrics"},
							{"containerPort": 9558},
							{"hostPort": 9558},
						}},
						{"volumeMounts": []*volumeMount{{
							Name:      "proc",
							MountPath: "/host/proc",
							ReadOnly:  true,
						}, {
							Name:      "systemd",
							MountPath: "/run/systemd",
							ReadOnly:  true,
						}}},
						{"resources": map[string]*resourceList{
							"limits": {
								Memory: "100Mi",
								Cpu:    "10m",
							},
							"requests": {
								Memory: "100Mi",
								Cpu:    "10m",
							},
						}},
					},
					"volumes": []*volume{{
						Name: "proc",
						HostPath: hostPath{
							Path: "/proc",
						},
					}, {
						Name: "systemd",
						HostPath: hostPath{
							Path: "/run/systemd",
						},
					}},
				},
			},
		},
	}
}

type volumeMount struct {
	Name      string
	MountPath string
	ReadOnly  bool
}

type resourceList struct {
	Memory string
	Cpu    string
}

type volume struct {
	Name     string
	HostPath hostPath
}

type hostPath struct {
	Path string
}
