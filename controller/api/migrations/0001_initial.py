# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations
import django_fsm
import api.models
import api.fields
from django.conf import settings
import json_field.fields


class Migration(migrations.Migration):

    dependencies = [
        migrations.swappable_dependency(settings.AUTH_USER_MODEL),
    ]

    operations = [
        migrations.CreateModel(
            name='App',
            fields=[
                ('uuid', api.fields.UuidField(auto_created=True, primary_key=True, serialize=False, editable=False, max_length=32, unique=True, verbose_name='UUID')),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('id', models.SlugField(unique=True, max_length=64)),
                ('structure', json_field.fields.JSONField(default={}, help_text='Enter a valid JSON object', blank=True, validators=[api.models.validate_app_structure])),
            ],
            options={
                'permissions': (('use_app', 'Can use app'),),
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Build',
            fields=[
                ('uuid', api.fields.UuidField(auto_created=True, primary_key=True, serialize=False, editable=False, max_length=32, unique=True, verbose_name='UUID')),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('image', models.CharField(max_length=256)),
                ('sha', models.CharField(max_length=40, blank=True)),
                ('procfile', json_field.fields.JSONField(default={}, help_text='Enter a valid JSON object', blank=True)),
                ('dockerfile', models.TextField(blank=True)),
                ('app', models.ForeignKey(to='api.App')),
                ('owner', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'ordering': ['-created'],
                'get_latest_by': 'created',
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Cluster',
            fields=[
                ('uuid', api.fields.UuidField(auto_created=True, primary_key=True, serialize=False, editable=False, max_length=32, unique=True, verbose_name='UUID')),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('id', models.CharField(unique=True, max_length=128)),
                ('type', models.CharField(default='coreos', max_length=16, choices=[('mock', 'Mock Cluster'), ('coreos', 'CoreOS Cluster'), ('faulty', 'Faulty Cluster')])),
                ('domain', models.CharField(max_length=128, validators=[api.models.validate_domain])),
                ('hosts', models.CharField(max_length=256, validators=[api.models.validate_comma_separated])),
                ('auth', models.TextField()),
                ('options', json_field.fields.JSONField(default={}, help_text='Enter a valid JSON object', blank=True)),
                ('owner', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'abstract': False,
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Config',
            fields=[
                ('uuid', api.fields.UuidField(auto_created=True, primary_key=True, serialize=False, editable=False, max_length=32, unique=True, verbose_name='UUID')),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('values', json_field.fields.JSONField(default={}, help_text='Enter a valid JSON object', blank=True)),
                ('memory', json_field.fields.JSONField(default={}, help_text='Enter a valid JSON object', blank=True)),
                ('cpu', json_field.fields.JSONField(default={}, help_text='Enter a valid JSON object', blank=True)),
                ('tags', json_field.fields.JSONField(default={}, help_text='Enter a valid JSON object', blank=True)),
                ('app', models.ForeignKey(to='api.App')),
                ('owner', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'ordering': ['-created'],
                'get_latest_by': 'created',
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Container',
            fields=[
                ('uuid', api.fields.UuidField(auto_created=True, primary_key=True, serialize=False, editable=False, max_length=32, unique=True, verbose_name='UUID')),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('type', models.CharField(max_length=128)),
                ('num', models.PositiveIntegerField()),
                ('state', django_fsm.FSMField(default='initialized', protected=True, max_length=50, choices=[('initialized', 'initialized'), ('created', 'created'), ('up', 'up'), ('down', 'down'), ('destroyed', 'destroyed')])),
                ('app', models.ForeignKey(to='api.App')),
                ('owner', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'ordering': ['created'],
                'get_latest_by': '-created',
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Domain',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('domain', models.TextField(unique=True)),
                ('app', models.ForeignKey(to='api.App')),
                ('owner', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'abstract': False,
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Key',
            fields=[
                ('uuid', api.fields.UuidField(auto_created=True, primary_key=True, serialize=False, editable=False, max_length=32, unique=True, verbose_name='UUID')),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('id', models.CharField(max_length=128)),
                ('public', models.TextField(unique=True)),
                ('owner', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'verbose_name': 'SSH Key',
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Push',
            fields=[
                ('uuid', api.fields.UuidField(auto_created=True, primary_key=True, serialize=False, editable=False, max_length=32, unique=True, verbose_name='UUID')),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('sha', models.CharField(max_length=40)),
                ('fingerprint', models.CharField(max_length=255)),
                ('receive_user', models.CharField(max_length=255)),
                ('receive_repo', models.CharField(max_length=255)),
                ('ssh_connection', models.CharField(max_length=255)),
                ('ssh_original_command', models.CharField(max_length=255)),
                ('app', models.ForeignKey(to='api.App')),
                ('owner', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'ordering': ['-created'],
                'get_latest_by': 'created',
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Release',
            fields=[
                ('uuid', api.fields.UuidField(auto_created=True, primary_key=True, serialize=False, editable=False, max_length=32, unique=True, verbose_name='UUID')),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('version', models.PositiveIntegerField()),
                ('summary', models.TextField(null=True, blank=True)),
                ('image', models.CharField(default='deis/helloworld', max_length=256)),
                ('app', models.ForeignKey(to='api.App')),
                ('build', models.ForeignKey(to='api.Build')),
                ('config', models.ForeignKey(to='api.Config')),
                ('owner', models.ForeignKey(to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'ordering': ['-created'],
                'get_latest_by': 'created',
            },
            bases=(models.Model,),
        ),
        migrations.AlterUniqueTogether(
            name='release',
            unique_together=set([('app', 'version')]),
        ),
        migrations.AlterUniqueTogether(
            name='push',
            unique_together=set([('app', 'uuid')]),
        ),
        migrations.AlterUniqueTogether(
            name='key',
            unique_together=set([('owner', 'id')]),
        ),
        migrations.AddField(
            model_name='container',
            name='release',
            field=models.ForeignKey(to='api.Release'),
            preserve_default=True,
        ),
        migrations.AlterUniqueTogether(
            name='config',
            unique_together=set([('app', 'uuid')]),
        ),
        migrations.AlterUniqueTogether(
            name='build',
            unique_together=set([('app', 'uuid')]),
        ),
        migrations.AddField(
            model_name='app',
            name='cluster',
            field=models.ForeignKey(to='api.Cluster'),
            preserve_default=True,
        ),
        migrations.AddField(
            model_name='app',
            name='owner',
            field=models.ForeignKey(to=settings.AUTH_USER_MODEL),
            preserve_default=True,
        ),
    ]
