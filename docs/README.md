# Deis Documentation

Source files for <http://docs.deis.io>, the documentation home of Deis.

Documentation in this tree consists of plain text files named with the
`.rst` suffix. These can be read in any text viewer, or processed by the
[Sphinx](http://sphinx-doc.org/) system to create an HTML version.

Please add any issues you find with this documentation to the
[Deis project](https://github.com/deis/deis/issues).

## Usage

1. Install Sphinx and other requirements:

    ```console
    $ virtualenv venv -q --prompt='(docs)' && . venv/bin/activate
    (docs)$ pip install -r docs_requirements.txt
    ```

2. Build the documentation and host it on a local web server:

    ```console
    (docs)$ make server
    sphinx-build -b dirhtml -d _build/doctrees   . _build/dirhtml
    Making output directory...
    Running Sphinx v1.3.0
    ...
    Build finished. The HTML pages are in _build/dirhtml.
    Serving HTTP on 0.0.0.0 port 8000 ...
    ```

3. Open a web browser to http://127.0.0.1:8000/ and learn about Deis,
a lightweight, flexible and powerful open source PaaS.

4. Fork this repository and send your changes or additions to Deis'
documentation as GitHub
[Pull Requests](https://github.com/deis/deis/pulls). Thank you!
