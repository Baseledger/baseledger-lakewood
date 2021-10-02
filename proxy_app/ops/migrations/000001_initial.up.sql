DO
$do$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE  rolname = 'baseledger') THEN
      CREATE ROLE baseledger WITH SUPERUSER LOGIN PASSWORD '<pass>';
    END IF;
END
$do$;

SET ROLE baseledger;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';

-- TODO: add PK and FK constraints

CREATE TABLE public.organizations (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  organization_name text NOT NULL
);

ALTER TABLE public.organizations OWNER TO baseledger;

ALTER TABLE ONLY public.organizations ADD CONSTRAINT organizations_pkey PRIMARY KEY (id);

INSERT INTO public.organizations (id, organization_name)
VALUES ('d45c9b93-3eef-4993-add6-aa1c84d17eea', 'Org1'), ('969e989c-bb61-4180-928c-0d48afd8c6a3', 'Org2'), ('4e227c7c-e73a-4e46-8907-786832cd87d3', 'Org3');

CREATE TABLE public.workgroups (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  workgroup_name text NOT NULL,
  privatize_key text NOT NULL
);

ALTER TABLE public.workgroups OWNER TO baseledger;

ALTER TABLE ONLY public.workgroups ADD CONSTRAINT workgroups_pkey PRIMARY KEY (id);

INSERT INTO public.workgroups (id, workgroup_name, privatize_key)
VALUES ('734276bc-4adc-4621-acf8-ac66dc91cb27', 'Workgroup1', '0c2e08bc9249fb42568e5a478e9af87a208471c46211a08f3ad9f0c5dbf57314');

CREATE TABLE public.workgroup_members (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  workgroup_id uuid NOT NULL,
  organization_id uuid NOT NULL,
  organization_endpoint text NOT NULL,
  organization_token text NOT NULL
);

ALTER TABLE public.workgroup_members OWNER TO baseledger;

ALTER TABLE ONLY public.workgroup_members ADD CONSTRAINT workgroup_members_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.workgroup_members
  ADD CONSTRAINT workgroup_members_organization_id_organizations_id_foreign FOREIGN KEY (organization_id) REFERENCES public.organizations(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE ONLY public.workgroup_members
  ADD CONSTRAINT workgroup_members_workgroup_id_workgroups_id_foreign FOREIGN KEY (workgroup_id) REFERENCES public.workgroups(id) ON UPDATE CASCADE ON DELETE CASCADE;

INSERT INTO public.workgroup_members (workgroup_id, organization_id, organization_endpoint, organization_token)
WITH
  w AS (
    SELECT id 
    FROM public.workgroups 
  ),
  o AS (
    SELECT id FROM public.organizations
    WHERE organization_name = 'Org1'
  )
  select w.id, o.id, 'host.docker.internal:4222', 'testToken1'
  from w, o;

INSERT INTO public.workgroup_members (workgroup_id, organization_id, organization_endpoint, organization_token)
WITH
  w AS (
    SELECT id 
    FROM public.workgroups 
  ),
  o AS (
    SELECT id FROM public.organizations
    WHERE organization_name = 'Org2'
  )
  select w.id, o.id, 'host.docker.internal:4223', 'testToken1'
  from w, o;

CREATE TABLE public.trustmeshes (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  created_at timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE public.trustmeshes OWNER TO baseledger;
ALTER TABLE ONLY public.trustmeshes ADD CONSTRAINT trustmeshes_pkey PRIMARY KEY (id);

CREATE TABLE public.offchain_process_messages (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  sender_id uuid,
	receiver_id uuid,
	topic text,
	referenced_offchain_process_message_id uuid,
	baseledger_sync_tree_json text,
	workstep_type text,
	business_object_proof text,
	tendermint_transaction_id_of_stored_proof uuid,
	baseledger_transaction_id_of_stored_proof uuid,
  baseledger_business_object_id text,
	referenced_baseledger_business_object_id text,
	status_text_message text,
  business_object_type text,
	baseledger_transaction_type text,
	referenced_baseledger_transaction_id uuid,
	entry_type text,
  sor_business_object_id text
);

ALTER TABLE public.offchain_process_messages OWNER TO baseledger;

ALTER TABLE ONLY public.offchain_process_messages ADD CONSTRAINT offchain_process_messages_pkey PRIMARY KEY (id);

CREATE TABLE public.trustmesh_entries (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
  tendermint_block_id text,
  tendermint_transaction_id uuid,
  tendermint_transaction_timestamp timestamp with time zone,
  
  entry_type text,

  sender_org_id uuid,
  receiver_org_id uuid,
  workgroup_id uuid,

  workstep_type text,
  baseledger_transaction_type text,

  baseledger_transaction_id uuid,
  referenced_baseledger_transaction_id uuid,

  business_object_type text,
  baseledger_business_object_id text,
  referenced_baseledger_business_object_id text,

  offchain_process_message_id uuid,
  referenced_process_message_id uuid,

  commitment_state text,
  transaction_hash text,
  trustmesh_id uuid NOT NULL,
  sor_business_object_id text
);

ALTER TABLE public.trustmesh_entries OWNER TO baseledger;

CREATE INDEX idx_commitment_state ON public.trustmesh_entries USING btree (commitment_state);

ALTER TABLE ONLY public.trustmesh_entries ADD CONSTRAINT trustmesh_entries_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.trustmesh_entries
  ADD CONSTRAINT trustmesh_entries_trustmesh_id_trustmeshes_id_foreign FOREIGN KEY (trustmesh_id) REFERENCES public.trustmeshes(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE ONLY public.trustmesh_entries
  ADD CONSTRAINT trustmesh_entries_workgroup_workgroups_id_foreign FOREIGN KEY (workgroup_id) REFERENCES public.workgroups(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE ONLY public.trustmesh_entries
  ADD CONSTRAINT trustmesh_entries_sender_org_organizations_id_foreign FOREIGN KEY (sender_org_id) REFERENCES public.organizations(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE ONLY public.trustmesh_entries
  ADD CONSTRAINT trustmesh_entries_receiver_org_organizations_id_foreign FOREIGN KEY (receiver_org_id) REFERENCES public.organizations(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE ONLY public.trustmesh_entries
  ADD CONSTRAINT trustmesh_entries_offchain_process_message_id_offchain_process_messages_id_foreign FOREIGN KEY (offchain_process_message_id) REFERENCES public.offchain_process_messages(id) ON UPDATE CASCADE ON DELETE CASCADE;

CREATE OR REPLACE FUNCTION set_trustmesh_entry_group()
  RETURNS trigger AS
  $$
    DECLARE new_trustmesh_id uuid;
    BEGIN
      IF NEW.referenced_baseledger_transaction_id = uuid_nil() THEN
        INSERT INTO trustmeshes VALUES (DEFAULT, DEFAULT) RETURNING id INTO new_trustmesh_id;
      ELSE 
        SELECT trustmesh_id INTO new_trustmesh_id FROM trustmesh_entries WHERE baseledger_transaction_id = NEW.referenced_baseledger_transaction_id LIMIT 1;
      END IF;
      NEW.trustmesh_id := new_trustmesh_id;
      RETURN NEW;
    END;
  $$
LANGUAGE plpgsql;

CREATE TRIGGER trustmesh_entry_insert_trigger
  BEFORE INSERT
  ON trustmesh_entries
  FOR EACH ROW
  EXECUTE PROCEDURE set_trustmesh_entry_group();


-- -- Add trustmeshes for visualization testing purposes
-- DO $$
-- DECLARE workgroup_id uuid;
-- BEGIN
--   SELECT id FROM workgroups INTO workgroup_id;

--   INSERT INTO public.offchain_process_messages (id, sender_id, receiver_id, topic, referenced_offchain_process_message_id, baseledger_sync_tree_json, workstep_type, business_object_proof, tendermint_transaction_id_of_stored_proof, baseledger_transaction_id_of_stored_proof, baseledger_business_object_id, referenced_baseledger_business_object_id, status_text_message, business_object_type, baseledger_transaction_type, referenced_baseledger_transaction_id, entry_type) VALUES ('c13607a7-40e4-4236-8b9b-ef8354b4c605', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '00000000-0000-0000-0000-000000000000', '{"RootProof":"1470c8c9e6f3a6c8e9ff8c1abb24edc9","Nodes":[{"SyncTreeNodeID":"478a3055-8b9e-4605-a1bf-9f0f85174de6","ParentNodeID":"99c098a6-5a82-458f-bc15-86c051f52cc1","Value":"Currency:USD","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":0},{"SyncTreeNodeID":"6fe46ef3-92c7-448d-9e37-28e8668c081d","ParentNodeID":"99c098a6-5a82-458f-bc15-86c051f52cc1","Value":"PurchaseOrderId:123","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":1},{"SyncTreeNodeID":"64fae2f2-f412-40a7-a882-c0a46d5a8e8a","ParentNodeID":"61366708-099a-4c00-b448-db57de63c79f","Value":"Amount:300","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":2},{"SyncTreeNodeID":"e10d3d62-4fed-4d09-84f1-fc9dc77ec4de","ParentNodeID":"61366708-099a-4c00-b448-db57de63c79f","Value":"","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":3},{"SyncTreeNodeID":"99c098a6-5a82-458f-bc15-86c051f52cc1","ParentNodeID":"18f2767e-b370-4bdd-acd8-e573f9835cba","Value":"4282ebfab5d7e569ebdb516d6cf2b769","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":0},{"SyncTreeNodeID":"61366708-099a-4c00-b448-db57de63c79f","ParentNodeID":"18f2767e-b370-4bdd-acd8-e573f9835cba","Value":"0929e875ffaf2c30d7825a4f523eab63","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":1},{"SyncTreeNodeID":"18f2767e-b370-4bdd-acd8-e573f9835cba","ParentNodeID":"","Value":"1470c8c9e6f3a6c8e9ff8c1abb24edc9","IsLeaf":false,"IsRoot":true,"IsHash":true,"IsCovered":false,"Level":2,"Index":0}]}', 'Initial', '1470c8c9e6f3a6c8e9ff8c1abb24edc9', 'c879f53f-8427-4d2a-a91a-7c68d7142129', 'c879f53f-8427-4d2a-a91a-7c68d7142129', 'a8d58233-7b8d-404b-8f1b-d044239f55ca', '00000000-0000-0000-0000-000000000000', 'Initial suggested', 'PurchaseOrder', 'Suggest', '00000000-0000-0000-0000-000000000000', 'SuggestionSent');
--   INSERT INTO public.offchain_process_messages (id, sender_id, receiver_id, topic, referenced_offchain_process_message_id, baseledger_sync_tree_json, workstep_type, business_object_proof, tendermint_transaction_id_of_stored_proof, baseledger_transaction_id_of_stored_proof, baseledger_business_object_id, referenced_baseledger_business_object_id, status_text_message, business_object_type, baseledger_transaction_type, referenced_baseledger_transaction_id, entry_type) VALUES ('9eb9c6b9-fea2-46eb-813e-30bf076176b7', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '68f6fb46-7fe5-4536-ac98-52b475418f7e', 'c13607a7-40e4-4236-8b9b-ef8354b4c605', '{"RootProof":"1470c8c9e6f3a6c8e9ff8c1abb24edc9","Nodes":[{"SyncTreeNodeID":"478a3055-8b9e-4605-a1bf-9f0f85174de6","ParentNodeID":"99c098a6-5a82-458f-bc15-86c051f52cc1","Value":"Currency:USD","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":0},{"SyncTreeNodeID":"6fe46ef3-92c7-448d-9e37-28e8668c081d","ParentNodeID":"99c098a6-5a82-458f-bc15-86c051f52cc1","Value":"PurchaseOrderId:123","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":1},{"SyncTreeNodeID":"64fae2f2-f412-40a7-a882-c0a46d5a8e8a","ParentNodeID":"61366708-099a-4c00-b448-db57de63c79f","Value":"Amount:300","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":2},{"SyncTreeNodeID":"e10d3d62-4fed-4d09-84f1-fc9dc77ec4de","ParentNodeID":"61366708-099a-4c00-b448-db57de63c79f","Value":"","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":3},{"SyncTreeNodeID":"99c098a6-5a82-458f-bc15-86c051f52cc1","ParentNodeID":"18f2767e-b370-4bdd-acd8-e573f9835cba","Value":"4282ebfab5d7e569ebdb516d6cf2b769","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":0},{"SyncTreeNodeID":"61366708-099a-4c00-b448-db57de63c79f","ParentNodeID":"18f2767e-b370-4bdd-acd8-e573f9835cba","Value":"0929e875ffaf2c30d7825a4f523eab63","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":1},{"SyncTreeNodeID":"18f2767e-b370-4bdd-acd8-e573f9835cba","ParentNodeID":"","Value":"1470c8c9e6f3a6c8e9ff8c1abb24edc9","IsLeaf":false,"IsRoot":true,"IsHash":true,"IsCovered":false,"Level":2,"Index":0}]}', 'Feedback', '', 'f0d02c3f-c5db-41a4-8e86-00db4114e436', 'f0d02c3f-c5db-41a4-8e86-00db4114e436', '00000000-0000-0000-0000-000000000000', 'a8d58233-7b8d-404b-8f1b-d044239f55ca', '', 'PurchaseOrder', 'Reject', 'c879f53f-8427-4d2a-a91a-7c68d7142129', 'FeedbackSent');
--   INSERT INTO public.offchain_process_messages (id, sender_id, receiver_id, topic, referenced_offchain_process_message_id, baseledger_sync_tree_json, workstep_type, business_object_proof, tendermint_transaction_id_of_stored_proof, baseledger_transaction_id_of_stored_proof, baseledger_business_object_id, referenced_baseledger_business_object_id, status_text_message, business_object_type, baseledger_transaction_type, referenced_baseledger_transaction_id, entry_type) VALUES ('b3b1615f-df38-4ea7-b48c-35c37411284d', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '00000000-0000-0000-0000-000000000000', '{"RootProof":"20b1355968b2a10f0c9c41c4f457f303","Nodes":[{"SyncTreeNodeID":"afd6513c-ff76-4a5f-b349-8fc8e13d9bf1","ParentNodeID":"8a1b9b01-dbbf-4d6d-ac7e-bbf2c45f76ba","Value":"PurchaseOrderId:123","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":0},{"SyncTreeNodeID":"38f0684a-e261-42c7-a59e-8f25b3229108","ParentNodeID":"8a1b9b01-dbbf-4d6d-ac7e-bbf2c45f76ba","Value":"Amount:500","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":1},{"SyncTreeNodeID":"1a029a2d-453d-47d3-9303-fbc87b20bcb6","ParentNodeID":"ad36863f-9c9a-40de-9c6c-8f3ed0600862","Value":"Currency:USD","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":2},{"SyncTreeNodeID":"1af48e41-e4cd-4622-abd8-462bf90eea6b","ParentNodeID":"ad36863f-9c9a-40de-9c6c-8f3ed0600862","Value":"","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":3},{"SyncTreeNodeID":"8a1b9b01-dbbf-4d6d-ac7e-bbf2c45f76ba","ParentNodeID":"c43441fd-7ccf-44ec-85c8-da36c58107bf","Value":"d7b94e8cdb31822101de69382867b3ee","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":0},{"SyncTreeNodeID":"ad36863f-9c9a-40de-9c6c-8f3ed0600862","ParentNodeID":"c43441fd-7ccf-44ec-85c8-da36c58107bf","Value":"48a43ed387efa8fb54af331aaeb4dcfd","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":1},{"SyncTreeNodeID":"c43441fd-7ccf-44ec-85c8-da36c58107bf","ParentNodeID":"","Value":"20b1355968b2a10f0c9c41c4f457f303","IsLeaf":false,"IsRoot":true,"IsHash":true,"IsCovered":false,"Level":2,"Index":0}]}', 'NewVersion', '20b1355968b2a10f0c9c41c4f457f303', '80ad989f-a47b-434d-9663-19a90097830c', '80ad989f-a47b-434d-9663-19a90097830c', 'cad15496-5400-44b9-b0c7-7e8d1e0d0740', 'a8d58233-7b8d-404b-8f1b-d044239f55ca', 'NewVersion suggested', 'PurchaseOrder', 'Suggest', 'f0d02c3f-c5db-41a4-8e86-00db4114e436', 'SuggestionSent');
--   INSERT INTO public.offchain_process_messages (id, sender_id, receiver_id, topic, referenced_offchain_process_message_id, baseledger_sync_tree_json, workstep_type, business_object_proof, tendermint_transaction_id_of_stored_proof, baseledger_transaction_id_of_stored_proof, baseledger_business_object_id, referenced_baseledger_business_object_id, status_text_message, business_object_type, baseledger_transaction_type, referenced_baseledger_transaction_id, entry_type) VALUES ('0aac68ca-f909-4f49-a932-c67468a5e211', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '68f6fb46-7fe5-4536-ac98-52b475418f7e', 'b3b1615f-df38-4ea7-b48c-35c37411284d', '{"RootProof":"20b1355968b2a10f0c9c41c4f457f303","Nodes":[{"SyncTreeNodeID":"afd6513c-ff76-4a5f-b349-8fc8e13d9bf1","ParentNodeID":"8a1b9b01-dbbf-4d6d-ac7e-bbf2c45f76ba","Value":"PurchaseOrderId:123","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":0},{"SyncTreeNodeID":"38f0684a-e261-42c7-a59e-8f25b3229108","ParentNodeID":"8a1b9b01-dbbf-4d6d-ac7e-bbf2c45f76ba","Value":"Amount:500","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":1},{"SyncTreeNodeID":"1a029a2d-453d-47d3-9303-fbc87b20bcb6","ParentNodeID":"ad36863f-9c9a-40de-9c6c-8f3ed0600862","Value":"Currency:USD","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":2},{"SyncTreeNodeID":"1af48e41-e4cd-4622-abd8-462bf90eea6b","ParentNodeID":"ad36863f-9c9a-40de-9c6c-8f3ed0600862","Value":"","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":3},{"SyncTreeNodeID":"8a1b9b01-dbbf-4d6d-ac7e-bbf2c45f76ba","ParentNodeID":"c43441fd-7ccf-44ec-85c8-da36c58107bf","Value":"d7b94e8cdb31822101de69382867b3ee","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":0},{"SyncTreeNodeID":"ad36863f-9c9a-40de-9c6c-8f3ed0600862","ParentNodeID":"c43441fd-7ccf-44ec-85c8-da36c58107bf","Value":"48a43ed387efa8fb54af331aaeb4dcfd","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":1},{"SyncTreeNodeID":"c43441fd-7ccf-44ec-85c8-da36c58107bf","ParentNodeID":"","Value":"20b1355968b2a10f0c9c41c4f457f303","IsLeaf":false,"IsRoot":true,"IsHash":true,"IsCovered":false,"Level":2,"Index":0}]}', 'Feedback', '', '70341890-80be-4e7e-b79d-079f8432dd7e', '70341890-80be-4e7e-b79d-079f8432dd7e', '00000000-0000-0000-0000-000000000000', 'cad15496-5400-44b9-b0c7-7e8d1e0d0740', '', 'PurchaseOrder', 'Approve', '80ad989f-a47b-434d-9663-19a90097830c', 'FeedbackSent');
--   INSERT INTO public.offchain_process_messages (id, sender_id, receiver_id, topic, referenced_offchain_process_message_id, baseledger_sync_tree_json, workstep_type, business_object_proof, tendermint_transaction_id_of_stored_proof, baseledger_transaction_id_of_stored_proof, baseledger_business_object_id, referenced_baseledger_business_object_id, status_text_message, business_object_type, baseledger_transaction_type, referenced_baseledger_transaction_id, entry_type) VALUES ('f014a80c-68cf-4c44-b813-b8a795428cb8', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '00000000-0000-0000-0000-000000000000', '{"RootProof":"1bc6ff002c4494df3f175a1e47168b47","Nodes":[{"SyncTreeNodeID":"3c25f772-ece7-421b-9b1b-09b2d41604d2","ParentNodeID":"814dce03-26f4-44c9-aeb2-4afc7912d249","Value":"InvoiceId:123","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":0},{"SyncTreeNodeID":"fefb9cd9-ec4b-4460-a9ac-d27be616f623","ParentNodeID":"814dce03-26f4-44c9-aeb2-4afc7912d249","Value":"Amount:100","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":1},{"SyncTreeNodeID":"63393a80-765f-4d18-b84c-f4f170f2b3d9","ParentNodeID":"d73e0eb2-95c0-4698-acba-9c7a1a1f1f80","Value":"Currency:EUR","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":2},{"SyncTreeNodeID":"517cd15b-6c39-488d-a668-0f68a8cee1d3","ParentNodeID":"d73e0eb2-95c0-4698-acba-9c7a1a1f1f80","Value":"","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":3},{"SyncTreeNodeID":"814dce03-26f4-44c9-aeb2-4afc7912d249","ParentNodeID":"ee012e53-449f-401a-93d3-1bfa26cac529","Value":"90ec038629c5cbd935711039c6db4baf","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":0},{"SyncTreeNodeID":"d73e0eb2-95c0-4698-acba-9c7a1a1f1f80","ParentNodeID":"ee012e53-449f-401a-93d3-1bfa26cac529","Value":"0373a984b45ff077fd6e3b9295b69ef2","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":1},{"SyncTreeNodeID":"ee012e53-449f-401a-93d3-1bfa26cac529","ParentNodeID":"","Value":"1bc6ff002c4494df3f175a1e47168b47","IsLeaf":false,"IsRoot":true,"IsHash":true,"IsCovered":false,"Level":2,"Index":0}]}', 'NextWorkstep', '1bc6ff002c4494df3f175a1e47168b47', 'a7d5659d-78fe-4093-8dbb-13807b660d30', 'a7d5659d-78fe-4093-8dbb-13807b660d30', '7b9ed108-3b6a-43e3-81ca-ad218af1573c', 'cad15496-5400-44b9-b0c7-7e8d1e0d0740', 'NextWorkstep suggested', 'Invoice', 'Suggest', '70341890-80be-4e7e-b79d-079f8432dd7e', 'SuggestionSent');
--   INSERT INTO public.offchain_process_messages (id, sender_id, receiver_id, topic, referenced_offchain_process_message_id, baseledger_sync_tree_json, workstep_type, business_object_proof, tendermint_transaction_id_of_stored_proof, baseledger_transaction_id_of_stored_proof, baseledger_business_object_id, referenced_baseledger_business_object_id, status_text_message, business_object_type, baseledger_transaction_type, referenced_baseledger_transaction_id, entry_type) VALUES ('a6a80004-3242-4c2a-8638-44176c4b5a8a', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '68f6fb46-7fe5-4536-ac98-52b475418f7e', 'f014a80c-68cf-4c44-b813-b8a795428cb8', '{"RootProof":"1bc6ff002c4494df3f175a1e47168b47","Nodes":[{"SyncTreeNodeID":"3c25f772-ece7-421b-9b1b-09b2d41604d2","ParentNodeID":"814dce03-26f4-44c9-aeb2-4afc7912d249","Value":"InvoiceId:123","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":0},{"SyncTreeNodeID":"fefb9cd9-ec4b-4460-a9ac-d27be616f623","ParentNodeID":"814dce03-26f4-44c9-aeb2-4afc7912d249","Value":"Amount:100","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":1},{"SyncTreeNodeID":"63393a80-765f-4d18-b84c-f4f170f2b3d9","ParentNodeID":"d73e0eb2-95c0-4698-acba-9c7a1a1f1f80","Value":"Currency:EUR","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":2},{"SyncTreeNodeID":"517cd15b-6c39-488d-a668-0f68a8cee1d3","ParentNodeID":"d73e0eb2-95c0-4698-acba-9c7a1a1f1f80","Value":"","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":3},{"SyncTreeNodeID":"814dce03-26f4-44c9-aeb2-4afc7912d249","ParentNodeID":"ee012e53-449f-401a-93d3-1bfa26cac529","Value":"90ec038629c5cbd935711039c6db4baf","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":0},{"SyncTreeNodeID":"d73e0eb2-95c0-4698-acba-9c7a1a1f1f80","ParentNodeID":"ee012e53-449f-401a-93d3-1bfa26cac529","Value":"0373a984b45ff077fd6e3b9295b69ef2","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":1},{"SyncTreeNodeID":"ee012e53-449f-401a-93d3-1bfa26cac529","ParentNodeID":"","Value":"1bc6ff002c4494df3f175a1e47168b47","IsLeaf":false,"IsRoot":true,"IsHash":true,"IsCovered":false,"Level":2,"Index":0}]}', 'Feedback', '', '62ad875f-a429-4bda-a9b1-cf4abcb59a8a', '62ad875f-a429-4bda-a9b1-cf4abcb59a8a', '00000000-0000-0000-0000-000000000000', '7b9ed108-3b6a-43e3-81ca-ad218af1573c', '', 'Invoice', 'Reject', 'a7d5659d-78fe-4093-8dbb-13807b660d30', 'FeedbackSent');
--   INSERT INTO public.offchain_process_messages (id, sender_id, receiver_id, topic, referenced_offchain_process_message_id, baseledger_sync_tree_json, workstep_type, business_object_proof, tendermint_transaction_id_of_stored_proof, baseledger_transaction_id_of_stored_proof, baseledger_business_object_id, referenced_baseledger_business_object_id, status_text_message, business_object_type, baseledger_transaction_type, referenced_baseledger_transaction_id, entry_type) VALUES ('eb5dbe73-64c6-4b84-92bc-ddf3683144b9', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '00000000-0000-0000-0000-000000000000', '{"RootProof":"60b52e956e622de329499bd4e9396424","Nodes":[{"SyncTreeNodeID":"cde6d7fb-e9a1-4c44-923b-00c084a790ac","ParentNodeID":"549db3d8-0509-4672-8254-bcdf3bb63c41","Value":"InvoiceId:123","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":0},{"SyncTreeNodeID":"6c8071aa-ec03-4261-92d5-1017e8f7c3f3","ParentNodeID":"549db3d8-0509-4672-8254-bcdf3bb63c41","Value":"Amount:300","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":1},{"SyncTreeNodeID":"8189e6d4-6631-4a93-97a4-8871c7c4183d","ParentNodeID":"71026ecd-7db0-4a69-84db-753c52888029","Value":"Currency:EUR","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":2},{"SyncTreeNodeID":"c36d1f0f-0ac7-4cd9-bdee-f6b666580cd2","ParentNodeID":"71026ecd-7db0-4a69-84db-753c52888029","Value":"","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":3},{"SyncTreeNodeID":"549db3d8-0509-4672-8254-bcdf3bb63c41","ParentNodeID":"08e645bf-fbac-4599-8b52-65a877b66db1","Value":"559edae586ae641cc533c1d40f834b22","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":0},{"SyncTreeNodeID":"71026ecd-7db0-4a69-84db-753c52888029","ParentNodeID":"08e645bf-fbac-4599-8b52-65a877b66db1","Value":"0373a984b45ff077fd6e3b9295b69ef2","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":1},{"SyncTreeNodeID":"08e645bf-fbac-4599-8b52-65a877b66db1","ParentNodeID":"","Value":"60b52e956e622de329499bd4e9396424","IsLeaf":false,"IsRoot":true,"IsHash":true,"IsCovered":false,"Level":2,"Index":0}]}', 'NewVersion', '60b52e956e622de329499bd4e9396424', '3596c0a9-d91e-49d8-9326-1c60d3bb01cc', '3596c0a9-d91e-49d8-9326-1c60d3bb01cc', '6a954c56-e85c-4d2b-bcd6-6d8b73d3797d', '7b9ed108-3b6a-43e3-81ca-ad218af1573c', 'NewVersion suggested', 'Invoice', 'Suggest', '62ad875f-a429-4bda-a9b1-cf4abcb59a8a', 'SuggestionSent');
--   INSERT INTO public.offchain_process_messages (id, sender_id, receiver_id, topic, referenced_offchain_process_message_id, baseledger_sync_tree_json, workstep_type, business_object_proof, tendermint_transaction_id_of_stored_proof, baseledger_transaction_id_of_stored_proof, baseledger_business_object_id, referenced_baseledger_business_object_id, status_text_message, business_object_type, baseledger_transaction_type, referenced_baseledger_transaction_id, entry_type) VALUES ('51d11026-ad6f-49a4-9870-db52f7d9d2db', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '68f6fb46-7fe5-4536-ac98-52b475418f7e', 'eb5dbe73-64c6-4b84-92bc-ddf3683144b9', '{"RootProof":"60b52e956e622de329499bd4e9396424","Nodes":[{"SyncTreeNodeID":"cde6d7fb-e9a1-4c44-923b-00c084a790ac","ParentNodeID":"549db3d8-0509-4672-8254-bcdf3bb63c41","Value":"InvoiceId:123","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":0},{"SyncTreeNodeID":"6c8071aa-ec03-4261-92d5-1017e8f7c3f3","ParentNodeID":"549db3d8-0509-4672-8254-bcdf3bb63c41","Value":"Amount:300","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":1},{"SyncTreeNodeID":"8189e6d4-6631-4a93-97a4-8871c7c4183d","ParentNodeID":"71026ecd-7db0-4a69-84db-753c52888029","Value":"Currency:EUR","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":2},{"SyncTreeNodeID":"c36d1f0f-0ac7-4cd9-bdee-f6b666580cd2","ParentNodeID":"71026ecd-7db0-4a69-84db-753c52888029","Value":"","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":3},{"SyncTreeNodeID":"549db3d8-0509-4672-8254-bcdf3bb63c41","ParentNodeID":"08e645bf-fbac-4599-8b52-65a877b66db1","Value":"559edae586ae641cc533c1d40f834b22","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":0},{"SyncTreeNodeID":"71026ecd-7db0-4a69-84db-753c52888029","ParentNodeID":"08e645bf-fbac-4599-8b52-65a877b66db1","Value":"0373a984b45ff077fd6e3b9295b69ef2","IsLeaf":false,"IsRoot":false,"IsHash":true,"IsCovered":false,"Level":1,"Index":1},{"SyncTreeNodeID":"08e645bf-fbac-4599-8b52-65a877b66db1","ParentNodeID":"","Value":"60b52e956e622de329499bd4e9396424","IsLeaf":false,"IsRoot":true,"IsHash":true,"IsCovered":false,"Level":2,"Index":0}]}', 'Feedback', '', '741fdd17-3f12-43be-a0ba-c0431ebe36c4', '741fdd17-3f12-43be-a0ba-c0431ebe36c4', '00000000-0000-0000-0000-000000000000', '6a954c56-e85c-4d2b-bcd6-6d8b73d3797d', '', 'Invoice', 'Approve', '3596c0a9-d91e-49d8-9326-1c60d3bb01cc', 'FeedbackSent');
--   INSERT INTO public.offchain_process_messages (id, sender_id, receiver_id, topic, referenced_offchain_process_message_id, baseledger_sync_tree_json, workstep_type, business_object_proof, tendermint_transaction_id_of_stored_proof, baseledger_transaction_id_of_stored_proof, baseledger_business_object_id, referenced_baseledger_business_object_id, status_text_message, business_object_type, baseledger_transaction_type, referenced_baseledger_transaction_id, entry_type) VALUES ('b802201e-5ca6-4e06-8d0e-67f84ebf6f5b', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '00000000-0000-0000-0000-000000000000', '{"RootProof":"fcb305074baf8ec24ae56dd7b48c4cd4","Nodes":[{"SyncTreeNodeID":"b2e89710-2c6a-4111-b4d3-0bdf39ab8ea0","ParentNodeID":"6f499ed3-d8e6-42f8-ab10-87e13a89c326","Value":"MTId:123","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":0},{"SyncTreeNodeID":"26513d0a-d154-40e4-ae54-43a87cab2159","ParentNodeID":"6f499ed3-d8e6-42f8-ab10-87e13a89c326","Value":"","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":1},{"SyncTreeNodeID":"6f499ed3-d8e6-42f8-ab10-87e13a89c326","ParentNodeID":"","Value":"fcb305074baf8ec24ae56dd7b48c4cd4","IsLeaf":false,"IsRoot":true,"IsHash":true,"IsCovered":false,"Level":1,"Index":0}]}', 'FinalWorkstep', 'fcb305074baf8ec24ae56dd7b48c4cd4', 'f2ba7772-5411-44cc-b358-49e3d893609b', 'f2ba7772-5411-44cc-b358-49e3d893609b', 'a28a27ce-ea1f-4134-bdb5-9065913afbb6', '6a954c56-e85c-4d2b-bcd6-6d8b73d3797d', 'FinalWorkstep suggested', 'MerkleTreeForExiting', 'Suggest', '741fdd17-3f12-43be-a0ba-c0431ebe36c4', 'SuggestionSent');
--   INSERT INTO public.offchain_process_messages (id, sender_id, receiver_id, topic, referenced_offchain_process_message_id, baseledger_sync_tree_json, workstep_type, business_object_proof, tendermint_transaction_id_of_stored_proof, baseledger_transaction_id_of_stored_proof, baseledger_business_object_id, referenced_baseledger_business_object_id, status_text_message, business_object_type, baseledger_transaction_type, referenced_baseledger_transaction_id, entry_type) VALUES ('001d2062-048b-4da6-92fb-47da89298406', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '68f6fb46-7fe5-4536-ac98-52b475418f7e', 'b802201e-5ca6-4e06-8d0e-67f84ebf6f5b', '{"RootProof":"fcb305074baf8ec24ae56dd7b48c4cd4","Nodes":[{"SyncTreeNodeID":"b2e89710-2c6a-4111-b4d3-0bdf39ab8ea0","ParentNodeID":"6f499ed3-d8e6-42f8-ab10-87e13a89c326","Value":"MTId:123","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":0},{"SyncTreeNodeID":"26513d0a-d154-40e4-ae54-43a87cab2159","ParentNodeID":"6f499ed3-d8e6-42f8-ab10-87e13a89c326","Value":"","IsLeaf":true,"IsRoot":false,"IsHash":false,"IsCovered":false,"Level":0,"Index":1},{"SyncTreeNodeID":"6f499ed3-d8e6-42f8-ab10-87e13a89c326","ParentNodeID":"","Value":"fcb305074baf8ec24ae56dd7b48c4cd4","IsLeaf":false,"IsRoot":true,"IsHash":true,"IsCovered":false,"Level":1,"Index":0}]}', 'Feedback', '', '85abe571-6d86-469c-9abd-91e5c8dbd221', '85abe571-6d86-469c-9abd-91e5c8dbd221', '00000000-0000-0000-0000-000000000000', 'a28a27ce-ea1f-4134-bdb5-9065913afbb6', '', 'MerkleTreeForExiting', 'Approve', 'f2ba7772-5411-44cc-b358-49e3d893609b', 'FeedbackSent');

--   INSERT INTO public.trustmesh_entries (id, tendermint_block_id, tendermint_transaction_id, tendermint_transaction_timestamp, entry_type, sender_org_id, receiver_org_id, workgroup_id, workstep_type, baseledger_transaction_type, baseledger_transaction_id, referenced_baseledger_transaction_id, business_object_type, baseledger_business_object_id, referenced_baseledger_business_object_id, offchain_process_message_id, referenced_process_message_id, commitment_state, transaction_hash, trustmesh_id) VALUES ('b9570432-171b-4634-a2bb-5f055d8b301b', '2478', 'c879f53f-8427-4d2a-a91a-7c68d7142129', '2021-07-27 22:55:19.100311+00', 'SuggestionSent', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '61ded832-b7ca-4100-8bc1-fb0935ff4436', workgroup_id, 'Initial', 'Suggest', 'c879f53f-8427-4d2a-a91a-7c68d7142129', '00000000-0000-0000-0000-000000000000', 'PurchaseOrder', 'a8d58233-7b8d-404b-8f1b-d044239f55ca', '00000000-0000-0000-0000-000000000000', 'c13607a7-40e4-4236-8b9b-ef8354b4c605', '00000000-0000-0000-0000-000000000000', 'COMMITTED', '36ACDEE6A36928C41585FD164A4FAC4CEF55FDFC5F9AD94E4FABCDE43F4C19FF', 'cbae6b84-09f4-4437-9290-3c459a7cb234');
--   INSERT INTO public.trustmesh_entries (id, tendermint_block_id, tendermint_transaction_id, tendermint_transaction_timestamp, entry_type, sender_org_id, receiver_org_id, workgroup_id, workstep_type, baseledger_transaction_type, baseledger_transaction_id, referenced_baseledger_transaction_id, business_object_type, baseledger_business_object_id, referenced_baseledger_business_object_id, offchain_process_message_id, referenced_process_message_id, commitment_state, transaction_hash, trustmesh_id) VALUES ('25d806f8-116b-4217-868b-dd5d15d840af', '2486', 'f0d02c3f-c5db-41a4-8e86-00db4114e436', '2021-07-27 22:55:24.196494+00', 'FeedbackSent', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', workgroup_id, 'Feedback', 'Reject', 'f0d02c3f-c5db-41a4-8e86-00db4114e436', 'c879f53f-8427-4d2a-a91a-7c68d7142129', 'PurchaseOrder', '00000000-0000-0000-0000-000000000000', 'a8d58233-7b8d-404b-8f1b-d044239f55ca', '9eb9c6b9-fea2-46eb-813e-30bf076176b7', 'c13607a7-40e4-4236-8b9b-ef8354b4c605', 'COMMITTED', '370AEC97182A559CE1B3B44D13712F81A27C521B08001A82F380A5C49CBAADE9', 'cbae6b84-09f4-4437-9290-3c459a7cb234');
--   INSERT INTO public.trustmesh_entries (id, tendermint_block_id, tendermint_transaction_id, tendermint_transaction_timestamp, entry_type, sender_org_id, receiver_org_id, workgroup_id, workstep_type, baseledger_transaction_type, baseledger_transaction_id, referenced_baseledger_transaction_id, business_object_type, baseledger_business_object_id, referenced_baseledger_business_object_id, offchain_process_message_id, referenced_process_message_id, commitment_state, transaction_hash, trustmesh_id) VALUES ('6f019dea-3924-4402-ad06-514db376eed9', '2502', '80ad989f-a47b-434d-9663-19a90097830c', '2021-07-27 22:55:43.542064+00', 'SuggestionSent', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '61ded832-b7ca-4100-8bc1-fb0935ff4436', workgroup_id, 'NewVersion', 'Suggest', '80ad989f-a47b-434d-9663-19a90097830c', 'f0d02c3f-c5db-41a4-8e86-00db4114e436', 'PurchaseOrder', 'cad15496-5400-44b9-b0c7-7e8d1e0d0740', 'a8d58233-7b8d-404b-8f1b-d044239f55ca', 'b3b1615f-df38-4ea7-b48c-35c37411284d', '00000000-0000-0000-0000-000000000000', 'COMMITTED', '893F27834A366DC491D686B9DC545D98D2DF4AC9080E018AC6DFEECD1232AB0C', 'cbae6b84-09f4-4437-9290-3c459a7cb234');
--   INSERT INTO public.trustmesh_entries (id, tendermint_block_id, tendermint_transaction_id, tendermint_transaction_timestamp, entry_type, sender_org_id, receiver_org_id, workgroup_id, workstep_type, baseledger_transaction_type, baseledger_transaction_id, referenced_baseledger_transaction_id, business_object_type, baseledger_business_object_id, referenced_baseledger_business_object_id, offchain_process_message_id, referenced_process_message_id, commitment_state, transaction_hash, trustmesh_id) VALUES ('12d7f4cc-f461-49d2-8057-1ab2f7e94322', '2516', '70341890-80be-4e7e-b79d-079f8432dd7e', '2021-07-27 22:55:58.822966+00', 'FeedbackSent', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', workgroup_id, 'Feedback', 'Approve', '70341890-80be-4e7e-b79d-079f8432dd7e', '80ad989f-a47b-434d-9663-19a90097830c', 'PurchaseOrder', '00000000-0000-0000-0000-000000000000', 'cad15496-5400-44b9-b0c7-7e8d1e0d0740', '0aac68ca-f909-4f49-a932-c67468a5e211', 'b3b1615f-df38-4ea7-b48c-35c37411284d', 'COMMITTED', '81557C1BCBBAC249B9FF55B2BDF861BC3F47C120877D54A34B04BEB8EC64A34C', 'cbae6b84-09f4-4437-9290-3c459a7cb234');
--   INSERT INTO public.trustmesh_entries (id, tendermint_block_id, tendermint_transaction_id, tendermint_transaction_timestamp, entry_type, sender_org_id, receiver_org_id, workgroup_id, workstep_type, baseledger_transaction_type, baseledger_transaction_id, referenced_baseledger_transaction_id, business_object_type, baseledger_business_object_id, referenced_baseledger_business_object_id, offchain_process_message_id, referenced_process_message_id, commitment_state, transaction_hash, trustmesh_id) VALUES ('6324194b-88d2-4dfd-89ac-c30f129c321d', '2545', 'a7d5659d-78fe-4093-8dbb-13807b660d30', '2021-07-27 22:56:24.283424+00', 'SuggestionSent', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '61ded832-b7ca-4100-8bc1-fb0935ff4436', workgroup_id, 'NextWorkstep', 'Suggest', 'a7d5659d-78fe-4093-8dbb-13807b660d30', '70341890-80be-4e7e-b79d-079f8432dd7e', 'Invoice', '7b9ed108-3b6a-43e3-81ca-ad218af1573c', 'cad15496-5400-44b9-b0c7-7e8d1e0d0740', 'f014a80c-68cf-4c44-b813-b8a795428cb8', '00000000-0000-0000-0000-000000000000', 'COMMITTED', '12E81C347F9C3D24ECF05DBEF202BFDFD2B410B510BD6440F180800DFED60BED', 'cbae6b84-09f4-4437-9290-3c459a7cb234');
--   INSERT INTO public.trustmesh_entries (id, tendermint_block_id, tendermint_transaction_id, tendermint_transaction_timestamp, entry_type, sender_org_id, receiver_org_id, workgroup_id, workstep_type, baseledger_transaction_type, baseledger_transaction_id, referenced_baseledger_transaction_id, business_object_type, baseledger_business_object_id, referenced_baseledger_business_object_id, offchain_process_message_id, referenced_process_message_id, commitment_state, transaction_hash, trustmesh_id) VALUES ('9365652f-231e-4eff-b029-961814141867', '2558', '62ad875f-a429-4bda-a9b1-cf4abcb59a8a', '2021-07-27 22:56:38.533982+00', 'FeedbackSent', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', workgroup_id, 'Feedback', 'Reject', '62ad875f-a429-4bda-a9b1-cf4abcb59a8a', 'a7d5659d-78fe-4093-8dbb-13807b660d30', 'Invoice', '00000000-0000-0000-0000-000000000000', '7b9ed108-3b6a-43e3-81ca-ad218af1573c', 'a6a80004-3242-4c2a-8638-44176c4b5a8a', 'f014a80c-68cf-4c44-b813-b8a795428cb8', 'COMMITTED', 'CFC3584CC24EE01C0D9C33C5CF9A6E0CB36E75185904448EB7AB237A16B786BB', 'cbae6b84-09f4-4437-9290-3c459a7cb234');
--   INSERT INTO public.trustmesh_entries (id, tendermint_block_id, tendermint_transaction_id, tendermint_transaction_timestamp, entry_type, sender_org_id, receiver_org_id, workgroup_id, workstep_type, baseledger_transaction_type, baseledger_transaction_id, referenced_baseledger_transaction_id, business_object_type, baseledger_business_object_id, referenced_baseledger_business_object_id, offchain_process_message_id, referenced_process_message_id, commitment_state, transaction_hash, trustmesh_id) VALUES ('98ea73ae-9899-4ee4-bb41-1e7ed11ed4e0', '2568', '3596c0a9-d91e-49d8-9326-1c60d3bb01cc', '2021-07-27 22:56:48.711176+00', 'SuggestionSent', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '61ded832-b7ca-4100-8bc1-fb0935ff4436', workgroup_id, 'NewVersion', 'Suggest', '3596c0a9-d91e-49d8-9326-1c60d3bb01cc', '62ad875f-a429-4bda-a9b1-cf4abcb59a8a', 'Invoice', '6a954c56-e85c-4d2b-bcd6-6d8b73d3797d', '7b9ed108-3b6a-43e3-81ca-ad218af1573c', 'eb5dbe73-64c6-4b84-92bc-ddf3683144b9', '00000000-0000-0000-0000-000000000000', 'COMMITTED', '16065CBEEC3B2D8C587FD4DB7AFE461EB32D28659ED1DD7F21E060FC2E47E7A0', 'cbae6b84-09f4-4437-9290-3c459a7cb234');
--   INSERT INTO public.trustmesh_entries (id, tendermint_block_id, tendermint_transaction_id, tendermint_transaction_timestamp, entry_type, sender_org_id, receiver_org_id, workgroup_id, workstep_type, baseledger_transaction_type, baseledger_transaction_id, referenced_baseledger_transaction_id, business_object_type, baseledger_business_object_id, referenced_baseledger_business_object_id, offchain_process_message_id, referenced_process_message_id, commitment_state, transaction_hash, trustmesh_id) VALUES ('215b4cd2-ae66-4a4e-beac-767fc01b9e6e', '2578', '741fdd17-3f12-43be-a0ba-c0431ebe36c4', '2021-07-27 22:56:58.889182+00', 'FeedbackSent', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', workgroup_id, 'Feedback', 'Approve', '741fdd17-3f12-43be-a0ba-c0431ebe36c4', '3596c0a9-d91e-49d8-9326-1c60d3bb01cc', 'Invoice', '00000000-0000-0000-0000-000000000000', '6a954c56-e85c-4d2b-bcd6-6d8b73d3797d', '51d11026-ad6f-49a4-9870-db52f7d9d2db', 'eb5dbe73-64c6-4b84-92bc-ddf3683144b9', 'COMMITTED', '9C2D79078A95E2B896BD654B4EF9AFDA76417733046E952BE9079DCF80FE23E3', 'cbae6b84-09f4-4437-9290-3c459a7cb234');
--   INSERT INTO public.trustmesh_entries (id, tendermint_block_id, tendermint_transaction_id, tendermint_transaction_timestamp, entry_type, sender_org_id, receiver_org_id, workgroup_id, workstep_type, baseledger_transaction_type, baseledger_transaction_id, referenced_baseledger_transaction_id, business_object_type, baseledger_business_object_id, referenced_baseledger_business_object_id, offchain_process_message_id, referenced_process_message_id, commitment_state, transaction_hash, trustmesh_id) VALUES ('ef014811-92ec-47a1-8c0e-9ba618ded03b', '2589', 'f2ba7772-5411-44cc-b358-49e3d893609b', '2021-07-27 22:57:09.071584+00', 'SuggestionSent', '68f6fb46-7fe5-4536-ac98-52b475418f7e', '61ded832-b7ca-4100-8bc1-fb0935ff4436', workgroup_id, 'FinalWorkstep', 'Suggest', 'f2ba7772-5411-44cc-b358-49e3d893609b', '741fdd17-3f12-43be-a0ba-c0431ebe36c4', 'MerkleTreeForExiting', 'a28a27ce-ea1f-4134-bdb5-9065913afbb6', '6a954c56-e85c-4d2b-bcd6-6d8b73d3797d', 'b802201e-5ca6-4e06-8d0e-67f84ebf6f5b', '00000000-0000-0000-0000-000000000000', 'COMMITTED', '61824D525D9290DE98C10E9CAA2A5BADA94B46E8F0B0B84F92480A528CD977FA', 'cbae6b84-09f4-4437-9290-3c459a7cb234');
--   INSERT INTO public.trustmesh_entries (id, tendermint_block_id, tendermint_transaction_id, tendermint_transaction_timestamp, entry_type, sender_org_id, receiver_org_id, workgroup_id, workstep_type, baseledger_transaction_type, baseledger_transaction_id, referenced_baseledger_transaction_id, business_object_type, baseledger_business_object_id, referenced_baseledger_business_object_id, offchain_process_message_id, referenced_process_message_id, commitment_state, transaction_hash, trustmesh_id) VALUES ('92dc60d9-400e-4639-bd19-2424550ef8ec', '2601', '85abe571-6d86-469c-9abd-91e5c8dbd221', '2021-07-27 22:57:24.336565+00', 'FeedbackSent', '61ded832-b7ca-4100-8bc1-fb0935ff4436', '68f6fb46-7fe5-4536-ac98-52b475418f7e', workgroup_id, 'Feedback', 'Approve', '85abe571-6d86-469c-9abd-91e5c8dbd221', 'f2ba7772-5411-44cc-b358-49e3d893609b', 'MerkleTreeForExiting', '00000000-0000-0000-0000-000000000000', 'a28a27ce-ea1f-4134-bdb5-9065913afbb6', '001d2062-048b-4da6-92fb-47da89298406', 'b802201e-5ca6-4e06-8d0e-67f84ebf6f5b', 'COMMITTED', '92281BE4575703AD85AB949FF9FFEB80293BF3C4B47952362AA50C42439BE472', 'cbae6b84-09f4-4437-9290-3c459a7cb234');

-- END $$;
