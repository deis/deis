deis-controller
===============

**Deis** is an open source *platform-as-a-service* (PaaS) for public and
private clouds.

Take your agile development to the next level. Free your mind and focus on
your code. Deploy updates to local metal or worldwide clouds with
```git push```. Scale servers, processes, and proxies with a simple command.
Enjoy the *twelve-factor app* workflow while keeping total control.

The [opdemand/deis-controller](https://github.com/opdemand/deis-controller)
project contains the RESTful API server. To set up your own private application
platform, the [deis-chef](https://github.com/opdemand/deis-chef) and
[deis](https://github.com/opdemand/deis) projects are also required.


Getting Started
---------------

First, make a clone of the
[deis-controller](https://github.com/opdemand/deis-controller)
github repository:

    git clone https://github.com/opdemand/deis-controller.git

This will create a **deis-controller** directory with the current code
and documentation. Change your current directory to the new project:

    cd deis-controller

Create a file for local application settings:

    touch deis/local_settings.py

Edit this new `deis/local_settings.py` file in a text editor, adding
the following section:

    # This is just an example for Sqlite3.
    DATABASES = {
	    'default': {
	        'ENGINE': 'django.db.backends.sqlite3',
	        'NAME': 'deis.db',
	        'USER': '',
	        'PASSWORD': '',
	        'HOST': '',
	        'PORT': '',
	    }
	}

Also add a value for SECRET_KEY:

	SECRET_KEY = 'atotallysecretkey'

Save these changes to deis/local_settings.py.

You might prefer to create a **virtualenv** to keep the project requirements
separate from other python installations. Assuming you have **virtualenv**
already installed:

    virtualenv venv --prompt='(deis)'
    source venv/bin/activate

from the deis project root will create a **virtualenv** in the *venv*
directory (which is ignored by `git` version control) and prepend the prompt
**(deis)** to the shell to remind you which environment you're in.

Whether or not you create a **virtualenv**, next use the `pip` tool to
install Django and other necessary python packages:

    pip install -U -r requirements.txt

Then create the database tables and indexes:

    python manage.py syncdb

Finally, run **Deis** in a test server:

    python manage.py runserver

You can simplify **Deis** development by making it your default
Django project. Set an environment variable in your .bashrc, .profile,
or the local equivalent:

    export DJANGO_SETTINGS_MODULE=deis.settings

Once you've done that, open a new command shell. You can create the
database tables and indexes with:

    make db

And you can run **Deis** in a test server with:

    make


License
-------

**Deis** is open source software under the Apache 2.0 license.
Please see the **LICENSE** file in the root directory for details.


Credits
-------

**Deis** rests on the shoulders of leading open source technologies:

  * Docker
  * Chef
  * Django
  * Heroku buildpacks
  * Gitosis

[OpDemand](http://www.opdemand.com/) sponsors and maintains the
**Deis** project.
