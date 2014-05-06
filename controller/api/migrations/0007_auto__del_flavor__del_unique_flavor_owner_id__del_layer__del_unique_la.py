# -*- coding: utf-8 -*-
from south.utils import datetime_utils as datetime
from south.db import db
from south.v2 import SchemaMigration
from django.db import models


class Migration(SchemaMigration):

    def forwards(self, orm):
        # Removing unique constraint on 'Container', fields ['formation', 'port']
        db.delete_unique(u'api_container', ['formation_id', 'port'])

        # Removing unique constraint on 'Container', fields ['app', 'type', 'num']
        db.delete_unique(u'api_container', ['app_id', 'type', 'num'])

        # Removing unique constraint on 'Config', fields ['app', 'version']
        db.delete_unique(u'api_config', ['app_id', 'version'])

        # Removing unique constraint on 'Formation', fields ['owner', 'id']
        db.delete_unique(u'api_formation', ['owner_id', 'id'])

        # Removing unique constraint on 'Node', fields ['formation', 'id']
        db.delete_unique(u'api_node', ['formation_id', 'id'])

        # Removing unique constraint on 'Provider', fields ['owner', 'id']
        db.delete_unique(u'api_provider', ['owner_id', 'id'])

        # Removing unique constraint on 'Layer', fields ['formation', 'id']
        db.delete_unique(u'api_layer', ['formation_id', 'id'])

        # Removing unique constraint on 'Flavor', fields ['owner', 'id']
        db.delete_unique(u'api_flavor', ['owner_id', 'id'])

        # Deleting model 'Flavor'
        db.delete_table(u'api_flavor')

        # Deleting model 'Layer'
        db.delete_table(u'api_layer')

        # Deleting model 'Provider'
        db.delete_table(u'api_provider')

        # Deleting model 'Node'
        db.delete_table(u'api_node')

        # Deleting model 'Formation'
        db.delete_table(u'api_formation')

        # Adding model 'Cluster'
        db.create_table(u'api_cluster', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('id', self.gf('django.db.models.fields.CharField')(unique=True, max_length=128)),
            ('type', self.gf('django.db.models.fields.CharField')(max_length=16)),
            ('domain', self.gf('django.db.models.fields.CharField')(max_length=128)),
            ('hosts', self.gf('django.db.models.fields.CharField')(max_length=256)),
            ('auth', self.gf('django.db.models.fields.TextField')()),
            ('options', self.gf('json_field.fields.JSONField')(default=u'{}', blank=True)),
        ))
        db.send_create_signal(u'api', ['Cluster'])


        # Changing field 'Release.build'
        db.alter_column(u'api_release', 'build_id', self.gf('django.db.models.fields.related.ForeignKey')(default='5e5dba0d-a7fe-4019-b392-62b7b993f1a8', to=orm['api.Build']))
        # Deleting field 'Config.version'
        db.delete_column(u'api_config', 'version')

        # Adding unique constraint on 'Config', fields ['app', 'uuid']
        db.create_unique(u'api_config', ['app_id', 'uuid'])

        # Deleting field 'Container.node'
        db.delete_column(u'api_container', 'node_id')

        # Deleting field 'Container.status'
        db.delete_column(u'api_container', 'status')

        # Deleting field 'Container.formation'
        db.delete_column(u'api_container', 'formation_id')

        # Deleting field 'Container.port'
        db.delete_column(u'api_container', 'port')

        # Adding field 'Container.release'
        db.add_column(u'api_container', 'release',
                      self.gf('django.db.models.fields.related.ForeignKey')(default='5e5dba0d-a7fe-4019-b392-62b7b993f1a8', to=orm['api.Release']),
                      keep_default=False)

        # Adding field 'Container.state'
        db.add_column(u'api_container', 'state',
                      self.gf('django.db.models.fields.CharField')(default=u'initializing', max_length=64),
                      keep_default=False)

        # Deleting field 'Build.procfile'
        db.delete_column(u'api_build', 'procfile')

        # Deleting field 'Build.size'
        db.delete_column(u'api_build', 'size')

        # Deleting field 'Build.url'
        db.delete_column(u'api_build', 'url')

        # Deleting field 'Build.checksum'
        db.delete_column(u'api_build', 'checksum')

        # Deleting field 'Build.dockerfile'
        db.delete_column(u'api_build', 'dockerfile')

        # Deleting field 'Build.sha'
        db.delete_column(u'api_build', 'sha')

        # Deleting field 'Build.output'
        db.delete_column(u'api_build', 'output')

        # Deleting field 'Build.config'
        db.delete_column(u'api_build', 'config')

        # Deleting field 'App.formation'
        db.delete_column(u'api_app', 'formation_id')

        # Deleting field 'App.containers'
        db.delete_column(u'api_app', 'containers')

        # Adding field 'App.cluster'
        db.add_column(u'api_app', 'cluster',
                      self.gf('django.db.models.fields.related.ForeignKey')(default='5e5dba0d-a7fe-4019-b392-62b7b993f1a8', to=orm['api.Cluster']),
                      keep_default=False)

        # Adding field 'App.structure'
        db.add_column(u'api_app', 'structure',
                      self.gf('json_field.fields.JSONField')(default=u'{}', blank=True),
                      keep_default=False)


    def backwards(self, orm):
        # Removing unique constraint on 'Config', fields ['app', 'uuid']
        db.delete_unique(u'api_config', ['app_id', 'uuid'])

        # Adding model 'Flavor'
        db.create_table(u'api_flavor', (
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('params', self.gf('json_field.fields.JSONField')(default=u'null', blank=True)),
            ('uuid', self.gf('api.fields.UuidField')(max_length=32, unique=True, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('provider', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Provider'])),
            ('id', self.gf('django.db.models.fields.SlugField')(max_length=64)),
        ))
        db.send_create_signal(u'api', ['Flavor'])

        # Adding unique constraint on 'Flavor', fields ['owner', 'id']
        db.create_unique(u'api_flavor', ['owner_id', 'id'])

        # Adding model 'Layer'
        db.create_table(u'api_layer', (
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('ssh_port', self.gf('django.db.models.fields.SmallIntegerField')(default=22)),
            ('ssh_username', self.gf('django.db.models.fields.CharField')(default=u'ubuntu', max_length=64)),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('proxy', self.gf('django.db.models.fields.BooleanField')(default=False)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('flavor', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Flavor'])),
            ('id', self.gf('django.db.models.fields.SlugField')(max_length=64)),
            ('uuid', self.gf('api.fields.UuidField')(max_length=32, unique=True, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('ssh_public_key', self.gf('django.db.models.fields.TextField')()),
            ('runtime', self.gf('django.db.models.fields.BooleanField')(default=False)),
            ('config', self.gf('json_field.fields.JSONField')(default=u'{}', blank=True)),
            ('ssh_private_key', self.gf('django.db.models.fields.TextField')()),
        ))
        db.send_create_signal(u'api', ['Layer'])

        # Adding unique constraint on 'Layer', fields ['formation', 'id']
        db.create_unique(u'api_layer', ['formation_id', 'id'])

        # Adding model 'Provider'
        db.create_table(u'api_provider', (
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('uuid', self.gf('api.fields.UuidField')(max_length=32, unique=True, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('type', self.gf('django.db.models.fields.SlugField')(max_length=16)),
            ('id', self.gf('django.db.models.fields.SlugField')(max_length=64)),
            ('creds', self.gf('json_field.fields.JSONField')(default=u'null', blank=True)),
        ))
        db.send_create_signal(u'api', ['Provider'])

        # Adding unique constraint on 'Provider', fields ['owner', 'id']
        db.create_unique(u'api_provider', ['owner_id', 'id'])

        # Adding model 'Node'
        db.create_table(u'api_node', (
            ('status', self.gf('json_field.fields.JSONField')(default=u'null', null=True, blank=True)),
            ('layer', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Layer'])),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('num', self.gf('django.db.models.fields.PositiveIntegerField')()),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('id', self.gf('django.db.models.fields.CharField')(max_length=64)),
            ('uuid', self.gf('api.fields.UuidField')(max_length=32, unique=True, primary_key=True)),
            ('provider_id', self.gf('django.db.models.fields.SlugField')(max_length=64, null=True, blank=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('fqdn', self.gf('django.db.models.fields.CharField')(max_length=256, null=True, blank=True)),
        ))
        db.send_create_signal(u'api', ['Node'])

        # Adding unique constraint on 'Node', fields ['formation', 'id']
        db.create_unique(u'api_node', ['formation_id', 'id'])

        # Adding model 'Formation'
        db.create_table(u'api_formation', (
            ('domain', self.gf('django.db.models.fields.CharField')(max_length=128, null=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('uuid', self.gf('api.fields.UuidField')(max_length=32, unique=True, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('nodes', self.gf('json_field.fields.JSONField')(default=u'{}', blank=True)),
            ('id', self.gf('django.db.models.fields.SlugField')(max_length=64, unique=True)),
        ))
        db.send_create_signal(u'api', ['Formation'])

        # Adding unique constraint on 'Formation', fields ['owner', 'id']
        db.create_unique(u'api_formation', ['owner_id', 'id'])

        # Deleting model 'Cluster'
        db.delete_table(u'api_cluster')


        # Changing field 'Release.build'
        db.alter_column(u'api_release', 'build_id', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Build'], null=True))

        # User chose to not deal with backwards NULL issues for 'Config.version'
        raise RuntimeError("Cannot reverse this migration. 'Config.version' and its values cannot be restored.")

        # The following code is provided here to aid in writing a correct migration        # Adding field 'Config.version'
        db.add_column(u'api_config', 'version',
                      self.gf('django.db.models.fields.PositiveIntegerField')(),
                      keep_default=False)

        # Adding unique constraint on 'Config', fields ['app', 'version']
        db.create_unique(u'api_config', ['app_id', 'version'])


        # User chose to not deal with backwards NULL issues for 'Container.node'
        raise RuntimeError("Cannot reverse this migration. 'Container.node' and its values cannot be restored.")

        # The following code is provided here to aid in writing a correct migration        # Adding field 'Container.node'
        db.add_column(u'api_container', 'node',
                      self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Node']),
                      keep_default=False)

        # Adding field 'Container.status'
        db.add_column(u'api_container', 'status',
                      self.gf('django.db.models.fields.CharField')(default=u'up', max_length=64),
                      keep_default=False)


        # User chose to not deal with backwards NULL issues for 'Container.formation'
        raise RuntimeError("Cannot reverse this migration. 'Container.formation' and its values cannot be restored.")

        # The following code is provided here to aid in writing a correct migration        # Adding field 'Container.formation'
        db.add_column(u'api_container', 'formation',
                      self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation']),
                      keep_default=False)


        # User chose to not deal with backwards NULL issues for 'Container.port'
        raise RuntimeError("Cannot reverse this migration. 'Container.port' and its values cannot be restored.")

        # The following code is provided here to aid in writing a correct migration        # Adding field 'Container.port'
        db.add_column(u'api_container', 'port',
                      self.gf('django.db.models.fields.PositiveIntegerField')(),
                      keep_default=False)

        # Deleting field 'Container.release'
        db.delete_column(u'api_container', 'release_id')

        # Deleting field 'Container.state'
        db.delete_column(u'api_container', 'state')

        # Adding unique constraint on 'Container', fields ['app', 'type', 'num']
        db.create_unique(u'api_container', ['app_id', 'type', 'num'])

        # Adding unique constraint on 'Container', fields ['formation', 'port']
        db.create_unique(u'api_container', ['formation_id', 'port'])

        # Adding field 'Build.procfile'
        db.add_column(u'api_build', 'procfile',
                      self.gf('json_field.fields.JSONField')(default=u'null', blank=True),
                      keep_default=False)

        # Adding field 'Build.size'
        db.add_column(u'api_build', 'size',
                      self.gf('django.db.models.fields.IntegerField')(null=True, blank=True),
                      keep_default=False)


        # User chose to not deal with backwards NULL issues for 'Build.url'
        raise RuntimeError("Cannot reverse this migration. 'Build.url' and its values cannot be restored.")

        # The following code is provided here to aid in writing a correct migration        # Adding field 'Build.url'
        db.add_column(u'api_build', 'url',
                      self.gf('django.db.models.fields.URLField')(max_length=200),
                      keep_default=False)

        # Adding field 'Build.checksum'
        db.add_column(u'api_build', 'checksum',
                      self.gf('django.db.models.fields.CharField')(default='', max_length=255, blank=True),
                      keep_default=False)

        # Adding field 'Build.dockerfile'
        db.add_column(u'api_build', 'dockerfile',
                      self.gf('django.db.models.fields.TextField')(default='', blank=True),
                      keep_default=False)

        # Adding field 'Build.sha'
        db.add_column(u'api_build', 'sha',
                      self.gf('django.db.models.fields.CharField')(default='', max_length=255, blank=True),
                      keep_default=False)

        # Adding field 'Build.output'
        db.add_column(u'api_build', 'output',
                      self.gf('django.db.models.fields.TextField')(default='', blank=True),
                      keep_default=False)

        # Adding field 'Build.config'
        db.add_column(u'api_build', 'config',
                      self.gf('json_field.fields.JSONField')(default=u'null', blank=True),
                      keep_default=False)


        # User chose to not deal with backwards NULL issues for 'App.formation'
        raise RuntimeError("Cannot reverse this migration. 'App.formation' and its values cannot be restored.")

        # The following code is provided here to aid in writing a correct migration        # Adding field 'App.formation'
        db.add_column(u'api_app', 'formation',
                      self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation']),
                      keep_default=False)

        # Adding field 'App.containers'
        db.add_column(u'api_app', 'containers',
                      self.gf('json_field.fields.JSONField')(default=u'{}', blank=True),
                      keep_default=False)

        # Deleting field 'App.cluster'
        db.delete_column(u'api_app', 'cluster_id')

        # Deleting field 'App.structure'
        db.delete_column(u'api_app', 'structure')


    models = {
        u'api.app': {
            'Meta': {'object_name': 'App'},
            'cluster': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Cluster']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'id': ('django.db.models.fields.SlugField', [], {'unique': 'True', 'max_length': '64'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'structure': ('json_field.fields.JSONField', [], {'default': "u'{}'", 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.build': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'app', u'uuid'),)", 'object_name': 'Build'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'image': ('django.db.models.fields.CharField', [], {'max_length': '256'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.cluster': {
            'Meta': {'object_name': 'Cluster'},
            'auth': ('django.db.models.fields.TextField', [], {}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'domain': ('django.db.models.fields.CharField', [], {'max_length': '128'}),
            'hosts': ('django.db.models.fields.CharField', [], {'max_length': '256'}),
            'id': ('django.db.models.fields.CharField', [], {'unique': 'True', 'max_length': '128'}),
            'options': ('json_field.fields.JSONField', [], {'default': "u'{}'", 'blank': 'True'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'type': ('django.db.models.fields.CharField', [], {'max_length': '16'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.config': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'app', u'uuid'),)", 'object_name': 'Config'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'}),
            'values': ('json_field.fields.JSONField', [], {'default': "u'{}'", 'blank': 'True'})
        },
        u'api.container': {
            'Meta': {'ordering': "[u'created']", 'object_name': 'Container'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'num': ('django.db.models.fields.PositiveIntegerField', [], {}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'release': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Release']"}),
            'state': ('django.db.models.fields.CharField', [], {'default': "u'initializing'", 'max_length': '64'}),
            'type': ('django.db.models.fields.CharField', [], {'max_length': '128', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.key': {
            'Meta': {'unique_together': "((u'owner', u'id'),)", 'object_name': 'Key'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'id': ('django.db.models.fields.CharField', [], {'max_length': '128'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'public': ('django.db.models.fields.TextField', [], {'unique': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.push': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'app', u'uuid'),)", 'object_name': 'Push'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'fingerprint': ('django.db.models.fields.CharField', [], {'max_length': '255'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'receive_repo': ('django.db.models.fields.CharField', [], {'max_length': '255'}),
            'receive_user': ('django.db.models.fields.CharField', [], {'max_length': '255'}),
            'sha': ('django.db.models.fields.CharField', [], {'max_length': '40'}),
            'ssh_connection': ('django.db.models.fields.CharField', [], {'max_length': '255'}),
            'ssh_original_command': ('django.db.models.fields.CharField', [], {'max_length': '255'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.release': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'app', u'version'),)", 'object_name': 'Release'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'build': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Build']"}),
            'config': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Config']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'summary': ('django.db.models.fields.TextField', [], {'null': 'True', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'}),
            'version': ('django.db.models.fields.PositiveIntegerField', [], {})
        },
        u'auth.group': {
            'Meta': {'object_name': 'Group'},
            u'id': ('django.db.models.fields.AutoField', [], {'primary_key': 'True'}),
            'name': ('django.db.models.fields.CharField', [], {'unique': 'True', 'max_length': '80'}),
            'permissions': ('django.db.models.fields.related.ManyToManyField', [], {'to': u"orm['auth.Permission']", 'symmetrical': 'False', 'blank': 'True'})
        },
        u'auth.permission': {
            'Meta': {'ordering': "(u'content_type__app_label', u'content_type__model', u'codename')", 'unique_together': "((u'content_type', u'codename'),)", 'object_name': 'Permission'},
            'codename': ('django.db.models.fields.CharField', [], {'max_length': '100'}),
            'content_type': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['contenttypes.ContentType']"}),
            u'id': ('django.db.models.fields.AutoField', [], {'primary_key': 'True'}),
            'name': ('django.db.models.fields.CharField', [], {'max_length': '50'})
        },
        u'auth.user': {
            'Meta': {'object_name': 'User'},
            'date_joined': ('django.db.models.fields.DateTimeField', [], {'default': 'datetime.datetime.now'}),
            'email': ('django.db.models.fields.EmailField', [], {'max_length': '75', 'blank': 'True'}),
            'first_name': ('django.db.models.fields.CharField', [], {'max_length': '30', 'blank': 'True'}),
            'groups': ('django.db.models.fields.related.ManyToManyField', [], {'symmetrical': 'False', 'related_name': "u'user_set'", 'blank': 'True', 'to': u"orm['auth.Group']"}),
            u'id': ('django.db.models.fields.AutoField', [], {'primary_key': 'True'}),
            'is_active': ('django.db.models.fields.BooleanField', [], {'default': 'True'}),
            'is_staff': ('django.db.models.fields.BooleanField', [], {'default': 'False'}),
            'is_superuser': ('django.db.models.fields.BooleanField', [], {'default': 'False'}),
            'last_login': ('django.db.models.fields.DateTimeField', [], {'default': 'datetime.datetime.now'}),
            'last_name': ('django.db.models.fields.CharField', [], {'max_length': '30', 'blank': 'True'}),
            'password': ('django.db.models.fields.CharField', [], {'max_length': '128'}),
            'user_permissions': ('django.db.models.fields.related.ManyToManyField', [], {'symmetrical': 'False', 'related_name': "u'user_set'", 'blank': 'True', 'to': u"orm['auth.Permission']"}),
            'username': ('django.db.models.fields.CharField', [], {'unique': 'True', 'max_length': '30'})
        },
        u'contenttypes.contenttype': {
            'Meta': {'ordering': "('name',)", 'unique_together': "(('app_label', 'model'),)", 'object_name': 'ContentType', 'db_table': "'django_content_type'"},
            'app_label': ('django.db.models.fields.CharField', [], {'max_length': '100'}),
            u'id': ('django.db.models.fields.AutoField', [], {'primary_key': 'True'}),
            'model': ('django.db.models.fields.CharField', [], {'max_length': '100'}),
            'name': ('django.db.models.fields.CharField', [], {'max_length': '100'})
        }
    }

    complete_apps = ['api']