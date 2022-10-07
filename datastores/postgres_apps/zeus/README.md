# Zeus Models #

High Level Explainer

### What is this repository for? ###

* Zeus Model Definitions and Interfaces

### Technology Stack ###

* Postgres 13.4 (Temporal Dependency)

### Directory
* Conversions
* Initializations
    * Seeds Useful Default Template Values
* Models
* Structs



* CI/CD
    * Flux

### Models

* Models

Model structs have two classes
* autogen
* manual

Use auto generation for creating base model structs using the codegen makefile

makefiles/codegen/Makefile

The manual structs should be higher level wrapper structs over the autogen ones

### Kubernetes -> SQL Guide

use naming to match k8s conventional naming: eg. StatefulSetSpec

synthetic chart subcomponent classes, eg container, ports, volume_mounts, etc
synthetic chart subcomponent classes are flexible, eg. ports as a component of a container
this allows hierarchy stacking

```text

chart_packages 
    one to many         -> chart_component_resources (chart component kind: e.g. service, statefulset, etc)
    which has many          -> chart_subcomponent_parent_class_types
    which has one to many       -> chart_subcomponent_child_class_types
    has at least one of             -> chart_subcomponents_child_values
                                    -> chart_subcomponents_jsonb_child_values
                                    -> chart_subcomponent_spec_pod_template_containers
                                    
Spec portions are just another parent class. PodSpecTemplates have elements added via jump tables

Deployment example hierarchy

              -> chart_subcomponent_parent_class_types (spec)
one to many      -> chart_subcomponent_child_class_types (deploymentSpec)                      
one to many        -> chart_subcomponents_child_values (volumes, nodeSelector, affinity, tolerations, etc)
zero to one        -> chart_subcomponent_spec_pod_template_containers (podTemplateSpec)
one to many        -> containers
one to many           -> containers_ports
one to many              -> container_ports
zero to many           -> containers_environmental_vars
one to many               -> container_environmental_vars
zero to many           -> containers_volume_mounts
one to many               -> container_volume_mounts
zero to many           -> containers_probes
one to many               -> container_probes
zero to one            -> container_compute_resources                   
```

****

Example
```yaml
spec:
 serviceName: "nginx"
 replicas: 2
 selector:
   matchLabels:
     app: nginx
 template: (chart_subcomponent_parent_class_types)
   metadata: (chart_subcomponent_child_class_types)
     labels: (chart_subcomponent_child_class_types)
       app: nginx (chart_subcomponents_child_values)
   spec: (chart_subcomponent_child_class_types)
     containers: (chart_subcomponent_child_class_types)
     - name: nginx (chart_subcomponents_child_values)
     image: registry.k8s.io/nginx-slim:0.8 (chart_subcomponents_child_values)
     ports: (chart_subcomponent_child_class_types)
      - containerPort: 80 (chart_subcomponents_child_values)
       name: web
     volumeMounts:
      - name: www
      mountPath: /usr/share/nginx/html

```
