package olympus_common_vals_cookbooks

import v1 "k8s.io/api/core/v1"

func GetChoreographyEnvVars() []v1.EnvVar {
	var envVars []v1.EnvVar

	refName := "choreography"

	name := "CLOUD_PROVIDER"
	key := "cloud-provider"
	envVars = append(envVars, MakeEnvVar(name, key, refName))

	name = "CTX"
	key = "ctx"
	envVars = append(envVars, MakeEnvVar(name, key, refName))

	name = "REGION"
	key = "region"
	envVars = append(envVars, MakeEnvVar(name, key, refName))

	name = "NS"
	key = "ns"
	envVars = append(envVars, MakeEnvVar(name, key, refName))
	return envVars
}

func GetCommonInternalAuthEnvVars() []v1.EnvVar {
	var envVars []v1.EnvVar

	name := "AGE_PKEY"
	key := "age-private-key"
	refName := "age-auth"
	envVars = append(envVars, MakeEnvVar(name, key, refName))

	name = "DO_SPACES_KEY"
	key = "do-spaces-key"
	refName = "spaces-key"
	envVars = append(envVars, MakeEnvVar(name, key, refName))

	name = "DO_SPACES_PKEY"
	key = "do-spaces-private-key"
	refName = "spaces-auth"
	envVars = append(envVars, MakeEnvVar(name, key, refName))

	return envVars
}

func MakeEnvVar(name, key, localObjRef string) v1.EnvVar {
	return v1.EnvVar{
		Name: name,
		ValueFrom: &v1.EnvVarSource{
			FieldRef:         nil,
			ResourceFieldRef: nil,
			ConfigMapKeyRef:  nil,
			SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{Name: localObjRef},
				Key:                  key,
			},
		},
	}
}
