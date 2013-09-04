# -*- coding: utf-8 -*-
import datetime
from south.db import db
from south.v2 import SchemaMigration
from django.db import models


class Migration(SchemaMigration):

    def forwards(self, orm):
        # Adding unique constraint on 'Key', fields ['public']
        db.create_unique(u'api_key', ['public'])


    def backwards(self, orm):
        # Removing unique constraint on 'Key', fields ['public']
        db.delete_unique(u'api_key', ['public'])


    models = {
        u'api.build': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'formation', u'uuid'),)", 'object_name': 'Build'},
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
            'url': ('django.db.models.fields.URLField', [], {'max_length': '200'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'})
        },
        u'api.config': {
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'formation', u'version'),)", 'object_name': 'Config'},
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'updated': ('django.db.models.fields.DateTimeField', [], {'auto_now': 'True', 'blank': 'True'}),
            'uuid': ('api.fields.UuidField', [], {'unique': 'True', 'max_length': '32', 'primary_key': 'True'}),
            'values': ('api.fields.EnvVarsField', [], {'default': "u'{}'", 'blank': 'True'}),
            'version': ('django.db.models.fields.PositiveIntegerField', [], {})
        },
        u'api.container': {
            'Meta': {'ordering': "[u'created']", 'unique_together': "((u'formation', u'type', u'num'),)", 'object_name': 'Container'},
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
            'containers': ('json_field.fields.JSONField', [], {'default': "u'{}'", 'blank': 'True'}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'id': ('django.db.models.fields.SlugField', [], {'max_length': '64'}),
            'layers': ('json_field.fields.JSONField', [], {'default': "u'{}'", 'blank': 'True'}),
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
            'chef_version': ('django.db.models.fields.CharField', [], {'default': "u'11.4.4'", 'max_length': '32'}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'environment': ('django.db.models.fields.CharField', [], {'default': "u'_default'", 'max_length': '64'}),
            'flavor': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Flavor']"}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
            'id': ('django.db.models.fields.SlugField', [], {'max_length': '64'}),
            'initial_attributes': ('json_field.fields.JSONField', [], {'default': "u'{}'", 'blank': 'True'}),
            'level': ('django.db.models.fields.PositiveIntegerField', [], {'default': '0'}),
            'owner': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['auth.User']"}),
            'run_list': ('django.db.models.fields.CharField', [], {'max_length': '512'}),
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
            'Meta': {'ordering': "[u'-created']", 'unique_together': "((u'formation', u'version'),)", 'object_name': 'Release'},
            'build': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Build']", 'null': 'True', 'blank': 'True'}),
            'config': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Config']"}),
            'created': ('django.db.models.fields.DateTimeField', [], {'auto_now_add': 'True', 'blank': 'True'}),
            'formation': ('django.db.models.fields.related.ForeignKey', [], {'to': u"orm['api.Formation']"}),
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