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
            ('public', self.gf('django.db.models.fields.TextField')()),
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
            ('creds', self.gf('api.fields.CredentialsField')(default=u'null')),
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
            ('ssh_username', self.gf('django.db.models.fields.CharField')(default=u'ubuntu', max_length=64)),
            ('ssh_private_key', self.gf('django.db.models.fields.TextField')()),
            ('ssh_public_key', self.gf('django.db.models.fields.TextField')()),
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
            ('flavor', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Flavor'])),
            ('image', self.gf('django.db.models.fields.CharField')(default=u'ubuntu', max_length=256)),
            ('structure', self.gf('json_field.fields.JSONField')(default=u'{}', blank=True)),
        ))
        db.send_create_signal(u'api', ['Formation'])

        # Adding unique constraint on 'Formation', fields ['owner', 'id']
        db.create_unique(u'api_formation', ['owner_id', 'id'])

        # Adding model 'Node'
        db.create_table(u'api_node', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('id', self.gf('django.db.models.fields.CharField')(max_length=64)),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('type', self.gf('django.db.models.fields.CharField')(default=u'backend', max_length=8)),
            ('num', self.gf('django.db.models.fields.PositiveIntegerField')(default=1)),
            ('provider_id', self.gf('django.db.models.fields.SlugField')(max_length=64, null=True, blank=True)),
            ('fqdn', self.gf('django.db.models.fields.CharField')(max_length=256, null=True, blank=True)),
            ('status', self.gf('api.fields.NodeStatusField')(default=u'null', null=True, blank=True)),
        ))
        db.send_create_signal(u'api', ['Node'])

        # Adding unique constraint on 'Node', fields ['owner', 'id']
        db.create_unique(u'api_node', ['owner_id', 'id'])

        # Adding model 'Backend'
        db.create_table(u'api_backend', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('node', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Node'])),
            ('status', self.gf('django.db.models.fields.CharField')(max_length=255, null=True, blank=True)),
        ))
        db.send_create_signal(u'api', ['Backend'])

        # Adding model 'Proxy'
        db.create_table(u'api_proxy', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('node', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Node'])),
            ('port', self.gf('django.db.models.fields.PositiveSmallIntegerField')()),
            ('in_proto', self.gf('django.db.models.fields.CharField')(default=u'HTTP', max_length=5)),
            ('out_proto', self.gf('django.db.models.fields.CharField')(default=u'HTTP', max_length=5)),
            ('flavor', self.gf('django.db.models.fields.CharField')(default=u'N', max_length=1)),
            ('status', self.gf('django.db.models.fields.CharField')(max_length=255, null=True, blank=True)),
        ))
        db.send_create_signal(u'api', ['Proxy'])

        # Adding model 'Container'
        db.create_table(u'api_container', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('node', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Node'])),
            ('type', self.gf('django.db.models.fields.CharField')(default=u'web', max_length=128)),
            ('num', self.gf('django.db.models.fields.PositiveIntegerField')()),
            ('id', self.gf('django.db.models.fields.CharField')(max_length=128, blank=True)),
            ('port', self.gf('django.db.models.fields.IntegerField')(null=True, blank=True)),
            ('metadata', self.gf('json_field.fields.JSONField')(default=u'null', blank=True)),
        ))
        db.send_create_signal(u'api', ['Container'])

        # Adding unique constraint on 'Container', fields ['formation', 'type', 'num']
        db.create_unique(u'api_container', ['formation_id', 'type', 'num'])

        # Adding model 'Config'
        db.create_table(u'api_config', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('version', self.gf('django.db.models.fields.PositiveIntegerField')(default=1)),
            ('values', self.gf('api.fields.EnvVarsField')(default=u'{}', blank=True)),
        ))
        db.send_create_signal(u'api', ['Config'])

        # Adding model 'Build'
        db.create_table(u'api_build', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('version', self.gf('django.db.models.fields.PositiveIntegerField')(default=1)),
            ('sha', self.gf('django.db.models.fields.CharField')(max_length=255, blank=True)),
            ('output', self.gf('django.db.models.fields.TextField')(blank=True)),
            ('procfile', self.gf('api.fields.ProcfileField')(default=u'null', blank=True)),
            ('dockerfile', self.gf('django.db.models.fields.TextField')(blank=True)),
            ('config', self.gf('api.fields.EnvVarsField')(default=u'null', blank=True)),
            ('url', self.gf('django.db.models.fields.URLField')(default=u'https://s3.amazonaws.com/gabrtv-slugs/nodejs.tar.gz', max_length=200)),
            ('size', self.gf('django.db.models.fields.IntegerField')(null=True, blank=True)),
            ('checksum', self.gf('django.db.models.fields.CharField')(max_length=255, blank=True)),
        ))
        db.send_create_signal(u'api', ['Build'])

        # Adding model 'Release'
        db.create_table(u'api_release', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('version', self.gf('django.db.models.fields.PositiveIntegerField')(default=1)),
            ('config', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Config'])),
            ('image', self.gf('django.db.models.fields.CharField')(default=u'ubuntu', max_length=256)),
            ('build', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Build'], null=True, blank=True)),
            ('args', self.gf('django.db.models.fields.CharField')(max_length=256, null=True, blank=True)),
            ('command', self.gf('django.db.models.fields.CharField')(max_length=256, null=True, blank=True)),
        ))
        db.send_create_signal(u'api', ['Release'])

        # Adding model 'Access'
        db.create_table(u'api_access', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
            ('role', self.gf('django.db.models.fields.CharField')(default=u'user', max_length=6)),
        ))
        db.send_create_signal(u'api', ['Access'])

        # Adding model 'Event'
        db.create_table(u'api_event', (
            ('uuid', self.gf('api.fields.UuidField')(unique=True, max_length=32, primary_key=True)),
            ('created', self.gf('django.db.models.fields.DateTimeField')(auto_now_add=True, blank=True)),
            ('updated', self.gf('django.db.models.fields.DateTimeField')(auto_now=True, blank=True)),
            ('owner', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['auth.User'])),
            ('formation', self.gf('django.db.models.fields.related.ForeignKey')(to=orm['api.Formation'])),
        ))
        db.send_create_signal(u'api', ['Event'])


    def backwards(self, orm):
        # Removing unique constraint on 'Container', fields ['formation', 'type', 'num']
        db.delete_unique(u'api_container', ['formation_id', 'type', 'num'])

        # Removing unique constraint on 'Node', fields ['owner', 'id']
        db.delete_unique(u'api_node', ['owner_id', 'id'])

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

        # Deleting model 'Node'
        db.delete_table(u'api_node')

        # Deleting model 'Backend'
        db.delete_table(u'api_backend')

        # Deleting model 'Proxy'
        db.delete_table(u'api_proxy')

        # Deleting model 'Container'
        db.delete_table(u'api_container')

        # Deleting model 'Config'
        db.delete_table(u'api_config')

        # Deleting model 'Build'
        db.delete_table(u'api_build')

        # Deleting model 'Release'
        db.delete_table(u'api_release')

        # Deleting model 'Access'
        db.delete_table(u'api_access')

        # Deleting model 'Event'
        db.delete_table(u'api_event')


    models = {
        u'api.access': {
            'Meta': {'object_name': 'Access'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'role': ('django.db.models.fields.CharField', [], {'default': "u'user'", 'max_length': '6'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.backend': {
            'Meta': {'object_name': 'Backend'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'node': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Node']"}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'status': ('django.db.models.fields.CharField', [], {'max_length': '255', 'null': 'True', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.build': {
            'Meta': {'ordering': "(u'-created',)", 'object_name': 'Build'},
            'checksum': ('django.db.models.fields.CharField', [], {'max_length': '255', 'blank': 'True'}),
            'config': ('api.fields.EnvVarsField', [], {'default': "u'null'", 'blank': 'True'}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'dockerfile': ('django.db.models.fields.TextField', [], {'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'output': ('django.db.models.fields.TextField', [], {'blank': 'True'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'procfile': ('api.fields.ProcfileField', [], {'default': "u'null'", 'blank': 'True'}),
            'sha': ('django.db.models.fields.CharField', [], {'max_length': '255', 'blank': 'True'}),
            'size': ('django.db.models.fields.IntegerField', [], {'null': 'True', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'url': ('django.db.models.fields.URLField', [], {'default': "u'https://s3.amazonaws.com/gabrtv-slugs/nodejs.tar.gz'", 'max_length': '200'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'}),
            'version': ('django.db.models.fields.PositiveIntegerField', [], {'default': '1'})
        },
        u'api.config': {
            'Meta': {'ordering': "[u'-created']", 'object_name': 'Config'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'}),
            'values': ('api.fields.EnvVarsField', [], {'default': "u'{}'", 'blank': 'True'}),
            'version': ('django.db.models.fields.PositiveIntegerField', [], {'default': '1'})
        },
        u'api.container': {
            'Meta': {'unique_together': "((u'formation', u'type', u'num'),)", 'object_name': 'Container'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'id': ('django.db.models.fields.CharField', [], {'max_length': '128', 'blank': 'True'}),
            'metadata': ('json_field.fields.JSONField', [], {'default': "u'null'", 'blank': 'True'}),
            'node': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Node']"}),
            'num': ('django.db.models.fields.PositiveIntegerField', [], {}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'port': ('django.db.models.fields.IntegerField', [], {'null': 'True', 'blank': 'True'}),
            'type': ('django.db.models.fields.CharField', [], {'default': "u'web'", 'max_length': '128'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.event': {
            'Meta': {'object_name': 'Event'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
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
            'ssh_private_key': ('django.db.models.fields.TextField', [], {}),
            'ssh_public_key': ('django.db.models.fields.TextField', [], {}),
            'ssh_username': ('django.db.models.fields.CharField', [], {'default': "u'ubuntu'", 'max_length': '64'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.formation': {
            'Meta': {'unique_together': "((u'owner', u'id'),)", 'object_name': 'Formation'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'flavor': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Flavor']"}),
            'id': ('django.db.models.fields.SlugField', [], {'max_length': '64'}),
            'image': ('django.db.models.fields.CharField', [], {'default': "u'ubuntu'", 'max_length': '256'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'structure': ('json_field.fields.JSONField', [], {'default': "u'{}'", 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.key': {
            'Meta': {'unique_together': "((u'owner', u'id'),)", 'object_name': 'Key'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'id': ('django.db.models.fields.CharField', [], {'max_length': '128'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'public': ('django.db.models.fields.TextField', [], {}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.node': {
            'Meta': {'unique_together': "((u'owner', u'id'),)", 'object_name': 'Node'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'fqdn': ('django.db.models.fields.CharField', [], {'max_length': '256', 'null': 'True', 'blank': 'True'}),
            'id': ('django.db.models.fields.CharField', [], {'max_length': '64'}),
            'num': ('django.db.models.fields.PositiveIntegerField', [], {'default': '1'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'provider_id': ('django.db.models.fields.SlugField', [], {'max_length': '64', 'null': 'True', 'blank': 'True'}),
            'status': ('api.fields.NodeStatusField', [], {'default': "u'null'", 'null': 'True', 'blank': 'True'}),
            'type': ('django.db.models.fields.CharField', [], {'default': "u'backend'", 'max_length': '8'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.provider': {
            'Meta': {'unique_together': "((u'owner', u'id'),)", 'object_name': 'Provider'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'creds': ('api.fields.CredentialsField', [], {'default': "u'null'"}),
            'id': ('django.db.models.fields.SlugField', [], {'max_length': '64'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'type': ('django.db.models.fields.SlugField', [], {'max_length': '16'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.proxy': {
            'Meta': {'object_name': 'Proxy'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'flavor': ('django.db.models.fields.CharField', [], {'default': "u'N'", 'max_length': '1'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'in_proto': ('django.db.models.fields.CharField', [], {'default': "u'HTTP'", 'max_length': '5'}),
            'node': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Node']"}),
            'out_proto': ('django.db.models.fields.CharField', [], {'default': "u'HTTP'", 'max_length': '5'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'port': ('django.db.models.fields.PositiveSmallIntegerField', [], {}),
            'status': ('django.db.models.fields.CharField', [], {'max_length': '255', 'null': 'True', 'blank': 'True'}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.release': {
            'Meta': {'ordering': "(u'-created',)", 'object_name': 'Release'},
            'args': ('django.db.models.fields.CharField', [], {'max_length': '256', 'null': 'True', 'blank': 'True'}),
            'build': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Build']", 'null': 'True', 'blank': 'True'}),
            'command': ('django.db.models.fields.CharField', [], {'max_length': '256', 'null': 'True', 'blank': 'True'}),
            'config': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Config']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'image': ('django.db.models.fields.CharField', [], {'default': "u'ubuntu'", 'max_length': '256'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'}),
            'version': ('django.db.models.fields.PositiveIntegerField', [], {'default': '1'})
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