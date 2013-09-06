# -*- coding: utf-8 -*-
import datetime
from south.db import db
from south.v2 import SchemaMigration
from django.db import models


class Migration(SchemaMigration):

    def forwards(self, orm):
        # Adding model 'Key'
        db.create_table(u'api_key', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('id', self.gf('django.db.models.fields.CharField')(max_length=128)),
            ('public', self.gf('django.db.models.fields.TextField')(unique=True)),
        ))
        db.send_create_signal(u'api', ['Key'])

        # Adding unique constraint on 'Key', fields ['owner', 'id']
        db.create_unique(u'api_key', ['owner_id', 'id'])

        # Adding model 'Provider'
        db.create_table(u'api_provider', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('id', self.gf('django.db.models.fields.SlugField')(max_length=64)),
            ('type', self.gf('django.db.models.fields.SlugField')(max_length=16)),
            ('creds', self.gf('api.fields.CredentialsField')(default=u'null', blank=True)),
        ))
        db.send_create_signal(u'api', ['Provider'])

        # Adding unique constraint on 'Provider', fields ['owner', 'id']
        db.create_unique(u'api_provider', ['owner_id', 'id'])

        # Adding model 'Flavor'
        db.create_table(u'api_flavor', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('id', self.gf('django.db.models.fields.SlugField')(max_length=64)),
            ('provider', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Provider'])),
            ('params', self.gf('api.fields.ParamsField')(default=u'null')),
            ('init', self.gf('api.fields.CloudInitField')()),
        ))
        db.send_create_signal(u'api', ['Flavor'])

        # Adding unique constraint on 'Flavor', fields ['owner', 'id']
        db.create_unique(u'api_flavor', ['owner_id', 'id'])

        # Adding model 'Formation'
        db.create_table(u'api_formation', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('id', self.gf('django.db.models.fields.SlugField')(max_length=64)),
            ('nodes', self.gf('json_field.fields.JSONField')(default=u'{}', blank=True)),
        ))
        db.send_create_signal(u'api', ['Formation'])

        # Adding unique constraint on 'Formation', fields ['owner', 'id']
        db.create_unique(u'api_formation', ['owner_id', 'id'])

        # Adding model 'Layer'
        db.create_table(u'api_layer', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('id', self.gf('django.db.models.fields.SlugField')(max_length=64)),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('flavor', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Flavor'])),
            ('proxy', self.gf('django.db.models.fields.BooleanField')(default=False)),
            ('runtime', self.gf('django.db.models.fields.BooleanField')(default=False)),
            ('ssh_username', self.gf('django.db.models.fields.CharField')(default=u'ubuntu', max_length=64)),
            ('ssh_private_key', self.gf('django.db.models.fields.TextField')()),
            ('ssh_public_key', self.gf('django.db.models.fields.TextField')()),
            ('config', self.gf('json_field.fields.JSONField')(default=u'{}', blank=True)),
        ))
        db.send_create_signal(u'api', ['Layer'])

        # Adding unique constraint on 'Layer', fields ['formation', 'id']
        db.create_unique(u'api_layer', ['formation_id', 'id'])

        # Adding model 'Node'
        db.create_table(u'api_node', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('id', self.gf('django.db.models.fields.CharField')(max_length=64)),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('layer', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Layer'])),
            ('num', self.gf('django.db.models.fields.PositiveIntegerField')()),
            ('provider_id', self.gf('django.db.models.fields.SlugField')(max_length=64, null=True, blank=True)),
            ('fqdn', self.gf('django.db.models.fields.CharField')(max_length=256, null=True, blank=True)),
            ('status', self.gf('api.fields.NodeStatusField')(default=u'null', null=True, blank=True)),
        ))
        db.send_create_signal(u'api', ['Node'])

        # Adding unique constraint on 'Node', fields ['formation', 'id']
        db.create_unique(u'api_node', ['formation_id', 'id'])

        # Adding model 'App'
        db.create_table(u'api_app', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('id', self.gf('django.db.models.fields.SlugField')(unique=True, max_length=64)),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('containers', self.gf('json_field.fields.JSONField')(default=u'{}', blank=True)),
        ))
        db.send_create_signal(u'api', ['App'])

        # Adding model 'Container'
        db.create_table(u'api_container', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('node', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Node'])),
            ('app', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.App'])),
            ('type', self.gf('django.db.models.fields.CharField')(max_length=128)),
            ('num', self.gf('django.db.models.fields.PositiveIntegerField')()),
            ('status', self.gf('django.db.models.fields.CharField')(default=u'up', max_length=64)),
        ))
        db.send_create_signal(u'api', ['Container'])

        # Adding unique constraint on 'Container', fields ['app', 'type', 'num']
        db.create_unique(u'api_container', ['app_id', 'type', 'num'])

        # Adding model 'Config'
        db.create_table(u'api_config', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('app', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.App'])),
            ('version', self.gf('django.db.models.fields.PositiveIntegerField')()),
            ('values', self.gf('api.fields.EnvVarsField')(default=u'{}', blank=True)),
        ))
        db.send_create_signal(u'api', ['Config'])

        # Adding unique constraint on 'Config', fields ['app', 'version']
        db.create_unique(u'api_config', ['app_id', 'version'])

        # Adding model 'Build'
        db.create_table(u'api_build', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('app', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.App'])),
            ('sha', self.gf('django.db.models.fields.CharField')(max_length=255, blank=True)),
            ('output', self.gf('django.db.models.fields.TextField')(blank=True)),
            ('procfile', self.gf('api.fields.ProcfileField')(default=u'null', blank=True)),
            ('dockerfile', self.gf('django.db.models.fields.TextField')(blank=True)),
            ('config', self.gf('api.fields.EnvVarsField')(default=u'null', blank=True)),
            ('url', self.gf('django.db.models.fields.URLField')(max_length=200)),
            ('size', self.gf('django.db.models.fields.IntegerField')(null=True, blank=True)),
            ('checksum', self.gf('django.db.models.fields.CharField')(max_length=255, blank=True)),
        ))
        db.send_create_signal(u'api', ['Build'])

        # Adding unique constraint on 'Build', fields ['app', 'uuid']
        db.create_unique(u'api_build', ['app_id', 'uuid'])

        # Adding model 'Release'
        db.create_table(u'api_release', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('app', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.App'])),
            ('version', self.gf('django.db.models.fields.PositiveIntegerField')()),
            ('config', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Config'])),
            ('image', self.gf('django.db.models.fields.CharField')(default=u'deis/buildstep', max_length=256)),
            ('build', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Build'], null=True, blank=True)),
        ))
        db.send_create_signal(u'api', ['Release'])

        # Adding unique constraint on 'Release', fields ['app', 'version']
        db.create_unique(u'api_release', ['app_id', 'version'])


    def backwards(self, orm):
        # Removing unique constraint on 'Release', fields ['app', 'version']
        db.delete_unique(u'api_release', ['app_id', 'version'])

        # Removing unique constraint on 'Build', fields ['app', 'uuid']
        db.delete_unique(u'api_build', ['app_id', 'uuid'])

        # Removing unique constraint on 'Config', fields ['app', 'version']
        db.delete_unique(u'api_config', ['app_id', 'version'])

        # Removing unique constraint on 'Container', fields ['app', 'type', 'num']
        db.delete_unique(u'api_container', ['app_id', 'type', 'num'])

        # Removing unique constraint on 'Node', fields ['formation', 'id']
        db.delete_unique(u'api_node', ['formation_id', 'id'])

        # Removing unique constraint on 'Layer', fields ['formation', 'id']
        db.delete_unique(u'api_layer', ['formation_id', 'id'])

        # Removing unique constraint on 'Formation', fields ['owner', 'id']
        db.delete_unique(u'api_formation', ['owner_id', 'id'])

        # Removing unique constraint on 'Flavor', fields ['owner', 'id']
        db.delete_unique(u'api_flavor', ['owner_id', 'id'])

        # Removing unique constraint on 'Provider', fields ['owner', 'id']
        db.delete_unique(u'api_provider', ['owner_id', 'id'])

        # Removing unique constraint on 'Key', fields ['owner', 'id']
        db.delete_unique(u'api_key', ['owner_id', 'id'])

        # Deleting model 'Key'
        db.delete_table(u'api_key')

        # Deleting model 'Provider'
        db.delete_table(u'api_provider')

        # Deleting model 'Flavor'
        db.delete_table(u'api_flavor')

        # Deleting model 'Formation'
        db.delete_table(u'api_formation')

        # Deleting model 'Layer'
        db.delete_table(u'api_layer')

        # Deleting model 'Node'
        db.delete_table(u'api_node')

        # Deleting model 'App'
        db.delete_table(u'api_app')

        # Deleting model 'Container'
        db.delete_table(u'api_container')

        # Deleting model 'Config'
        db.delete_table(u'api_config')

        # Deleting model 'Build'
        db.delete_table(u'api_build')

        # Deleting model 'Release'
        db.delete_table(u'api_release')


    models = {
        u'api.app': {
            'Meta': {'object_name': 'App'},
            'containers': ('json_field.fields.JSONField', [], {'default': "u'{}'", 'blank': 'True'}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'id': ('django.db.models.fields.SlugField', [], {'unique': 'True', 'max_length': '64'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.build': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'app', u'uuid'),)", 'object_name': 'Build'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'checksum': ('django.db.models.fields.CharField', [], {'max_length': '255', 'blank': 'True'}),
            'config': ('api.fields.EnvVarsField', [], {'default': "u'null'", 'blank': 'True'}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'dockerfile': ('django.db.models.fields.TextField', [], {'blank': 'True'}),
            'output': ('django.db.models.fields.TextField', [], {'blank': 'True'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'procfile': ('api.fields.ProcfileField', [], {'default': "u'null'", 'blank': 'True'}),
            'sha': ('django.db.models.fields.CharField', [], {'max_length': '255', 'blank': 'True'}),
            'size': ('django.db.models.fields.IntegerField', [], {'null': 'True', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'url': ('django.db.models.fields.URLField', [], {'max_length': '200'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.config': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'app', u'version'),)", 'object_name': 'Config'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'}),
            'values': ('api.fields.EnvVarsField', [], {'default': "u'{}'", 'blank': 'True'}),
            'version': ('django.db.models.fields.PositiveIntegerField', [], {})
        },
        u'api.container': {
            'Meta': {'ordering': "[u'created']", 'unique_together': "((u'app', u'type', u'num'),)", 'object_name': 'Container'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'node': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Node']"}),
            'num': ('django.db.models.fields.PositiveIntegerField', [], {}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'status': ('django.db.models.fields.CharField', [], {'default': "u'up'", 'max_length': '64'}),
            'type': ('django.db.models.fields.CharField', [], {'max_length': '128'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.flavor': {
            'Meta': {'unique_together': "((u'owner', u'id'),)", 'object_name': 'Flavor'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'id': ('django.db.models.fields.SlugField', [], {'max_length': '64'}),
            'init': ('api.fields.CloudInitField', [], {}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'params': ('api.fields.ParamsField', [], {'default': "u'null'"}),
            'provider': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Provider']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.formation': {
            'Meta': {'unique_together': "((u'owner', u'id'),)", 'object_name': 'Formation'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'id': ('django.db.models.fields.SlugField', [], {'max_length': '64'}),
            'nodes': ('json_field.fields.JSONField', [], {'default': "u'{}'", 'blank': 'True'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
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
        u'api.layer': {
            'Meta': {'unique_together': "((u'formation', u'id'),)", 'object_name': 'Layer'},
            'config': ('json_field.fields.JSONField', [], {'default': "u'{}'", 'blank': 'True'}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'flavor': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Flavor']"}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'id': ('django.db.models.fields.SlugField', [], {'max_length': '64'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'proxy': ('django.db.models.fields.BooleanField', [], {'default': 'False'}),
            'runtime': ('django.db.models.fields.BooleanField', [], {'default': 'False'}),
            'ssh_private_key': ('django.db.models.fields.TextField', [], {}),
            'ssh_public_key': ('django.db.models.fields.TextField', [], {}),
            'ssh_username': ('django.db.models.fields.CharField', [], {'default': "u'ubuntu'", 'max_length': '64'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.node': {
            'Meta': {'unique_together': "((u'formation', u'id'),)", 'object_name': 'Node'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'fqdn': ('django.db.models.fields.CharField', [], {'max_length': '256', 'null': 'True', 'blank': 'True'}),
            'id': ('django.db.models.fields.CharField', [], {'max_length': '64'}),
            'layer': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Layer']"}),
            'num': ('django.db.models.fields.PositiveIntegerField', [], {}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'provider_id': ('django.db.models.fields.SlugField', [], {'max_length': '64', 'null': 'True', 'blank': 'True'}),
            'status': ('api.fields.NodeStatusField', [], {'default': "u'null'", 'null': 'True', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.provider': {
            'Meta': {'unique_together': "((u'owner', u'id'),)", 'object_name': 'Provider'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'creds': ('api.fields.CredentialsField', [], {'default': "u'null'", 'blank': 'True'}),
            'id': ('django.db.models.fields.SlugField', [], {'max_length': '64'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'type': ('django.db.models.fields.SlugField', [], {'max_length': '16'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.release': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'app', u'version'),)", 'object_name': 'Release'},
            'app': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.App']"}),
            'build': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Build']", 'null': 'True', 'blank': 'True'}),
            'config': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Config']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'image': ('django.db.models.fields.CharField', [], {'default': "u'deis/buildstep'", 'max_length': '256'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
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
            'groups': ('django.db.models.fields.related.ManyToManyField', [], {'to': u"orm['auth.Group']", 'symmetrical': 'False', 'blank': 'True'}),
            u'id': ('django.db.models.fields.AutoField', [], {'primary_key': 'True'}),
            'is_active': ('django.db.models.fields.BooleanField', [], {'default': 'True'}),
            'is_staff': ('django.db.models.fields.BooleanField', [], {'default': 'False'}),
            'is_superuser': ('django.db.models.fields.BooleanField', [], {'default': 'False'}),
            'last_login': ('django.db.models.fields.DateTimeField', [], {'default': 'datetime.datetime.now'}),
            'last_name': ('django.db.models.fields.CharField', [], {'max_length': '30', 'blank': 'True'}),
            'password': ('django.db.models.fields.CharField', [], {'max_length': '128'}),
            'user_permissions': ('django.db.models.fields.related.ManyToManyField', [], {'to': u"orm['auth.Permission']", 'symmetrical': 'False', 'blank': 'True'}),
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