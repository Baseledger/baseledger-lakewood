DROP INDEX idx_commitment_state;


ALTER TABLE ONLY public.workgroup_members DROP CONSTRAINT workgroup_members_pkey;
ALTER TABLE ONLY public.trustmesh_entries DROP CONSTRAINT trustmesh_entries_workgroup_workgroups_id_foreign;
ALTER TABLE ONLY public.trustmesh_entries DROP CONSTRAINT trustmesh_entries_sender_org_organizations_id_foreign;
ALTER TABLE ONLY public.trustmesh_entries DROP CONSTRAINT trustmesh_entries_receiver_org_organizations_id_foreign;
ALTER TABLE ONLY public.workgroup_members DROP CONSTRAINT workgroup_members_organization_id_organizations_id_foreign;
ALTER TABLE ONLY public.workgroup_members DROP CONSTRAINT workgroup_members_workgroup_id_workgroups_id_foreign;
ALTER TABLE ONLY public.workgroups DROP CONSTRAINT workgroups_pkey;
ALTER TABLE ONLY public.organizations DROP CONSTRAINT organizations_pkey;

ALTER TABLE ONLY public.trustmesh_entries DROP CONSTRAINT trustmesh_entries_offchain_process_message_id_offchain_process_messages_id_foreign;
ALTER TABLE ONLY public.offchain_process_messages DROP CONSTRAINT offchain_process_messages_pkey;
ALTER TABLE ONLY public.trustmesh_entries DROP CONSTRAINT trustmesh_entries_pkey;
ALTER TABLE ONLY public.trustmesh_entries DROP CONSTRAINT trustmesh_entries_trustmesh_id_trustmeshes_id_foreign;
ALTER TABLE ONLY public.trustmeshes DROP CONSTRAINT trustmeshes_pkey;

DROP FUNCTION IF EXISTS set_trustmesh_entry_group CASCADE;

DROP TABLE public.workgroups;
DROP TABLE public.workgroup_members;
DROP TABLE public.trustmesh_entries;
DROP TABLE public.offchain_process_messages;
DROP TABLE public.organizations;
DROP TABLE public.trustmeshes;