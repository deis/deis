import os.path
import shutil
import subprocess
import tempfile


def publish_release(src_image, config, target_image):
    """
    Publish a new release as a Docker image

    Given a source image and dictionary of last-mile configuration,
    create a target Docker image on the registry.

    For example publish_release('registry.local:5000/gabrtv/myapp:<sha>',
                                {'ENVVAR': 'values'},
                                'registry.local:5000/gabrtv/myapp:v23',
    results in a new Docker image at 'registry.local:5000/gabrtv/myapp:v23' which
    contains the new configuration as ENV entries.
    """
    # write out dockerfile
    dockerfile = _build_dockerfile(src_image, config)
    tempdir = tempfile.mkdtemp()
    dockerfile_path = os.path.join(tempdir, 'Dockerfile')
    with open(dockerfile_path, 'w') as f:
        f.write(dockerfile)
    try:
        # pull the source image to ensure we have latest
        p = subprocess.Popen(['docker', 'pull', src_image])
        rc = p.wait()
        if rc != 0:
            raise RuntimeError('Failed to pull source image')
        # build the new image with last-mile configuration
        p = subprocess.Popen(['docker', 'build', '-t', target_image, tempdir])
        rc = p.wait()
        if rc != 0:
            raise RuntimeError('Failed to build release image')
        # push the target image
        p = subprocess.Popen(['docker', 'push', target_image])
        rc = p.wait()
        if rc != 0:
            raise RuntimeError('Failed to push release image')
    finally:
        shutil.rmtree(tempdir)
        # cleanup the temporary image
        p = subprocess.Popen(['docker', 'rmi', '-f', target_image])
        rc = p.wait()
        if rc != 0:
            print('warning: failed to delete temporary images')


def _build_dockerfile(image, config):
    dockerfile = ["FROM "+image]
    for k, v in config.items():
        dockerfile.append("ENV {} {}".format(k.upper(), v))
    return '\n'.join(dockerfile)
