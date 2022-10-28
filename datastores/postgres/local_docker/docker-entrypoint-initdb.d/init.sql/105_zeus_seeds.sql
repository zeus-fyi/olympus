BEGIN;
INSERT INTO "public"."chart_component_resources" VALUES (0, 'Deployment', 'apps/v1');
INSERT INTO "public"."chart_component_resources" VALUES (1, 'StatefulSet', 'apps/v1');
INSERT INTO "public"."chart_component_resources" VALUES (2, 'Service', 'v1');
INSERT INTO "public"."chart_component_resources" VALUES (3, 'ReplicaSet', 'apps/v1');
INSERT INTO "public"."chart_component_resources" VALUES (4, 'DaemonSet', 'apps/v1');
INSERT INTO "public"."chart_component_resources" VALUES (5, 'Job', 'batch/v1');
INSERT INTO "public"."chart_component_resources" VALUES (6, 'CronJob', 'batch/v1');
INSERT INTO "public"."chart_component_resources" VALUES (7, 'ReplicationController', 'v1');
INSERT INTO "public"."chart_component_resources" VALUES (8, 'PersistentVolumeClaim', 'v1');
INSERT INTO "public"."chart_component_resources" VALUES (9, 'PersistentVolume', 'v1');
INSERT INTO "public"."chart_component_resources" VALUES (10, 'StorageClass', 'storage.k8s.io/v1');
INSERT INTO "public"."chart_component_resources" VALUES (11, 'VolumeSnapshotContent', 'snapshot.storage.k8s.io/v1');
INSERT INTO "public"."chart_component_resources" VALUES (12, 'ConfigMap', 'v1');
INSERT INTO "public"."chart_component_resources" VALUES (13, 'Secret', 'v1');
INSERT INTO "public"."chart_component_resources" VALUES (14, 'Ingress', 'networking.k8s.io/v1');
INSERT INTO "public"."chart_component_resources" VALUES (15, 'IngressClass', 'networking.k8s.io/v1');
INSERT INTO "public"."chart_component_resources" VALUES (15, 'Pod', 'v1');
INSERT INTO "public"."chart_component_resources" VALUES (16, 'ServiceAccount', 'v1');
INSERT INTO "public"."chart_component_resources" VALUES (17, 'ClusterRole', 'rbac.authorization.k8s.io/v1');
INSERT INTO "public"."chart_component_resources" VALUES (18, 'RoleBinding', 'rbac.authorization.k8s.io/v1');
INSERT INTO "public"."chart_component_resources" VALUES (19, 'ClusterRoleBinding', 'rbac.authorization.k8s.io/v1');
INSERT INTO "public"."chart_component_resources" VALUES (20, 'Role', 'rbac.authorization.k8s.io/v1');
INSERT INTO "public"."chart_component_resources" VALUES (21, 'CSIDriver', 'storage.k8s.io/v1');
INSERT INTO "public"."chart_component_resources" VALUES (22, 'CSINode', 'storage.k8s.io/v1');
INSERT INTO "public"."chart_component_resources" VALUES (23, 'CSIStorageCapacity', 'storage.k8s.io/v1');


COMMIT;