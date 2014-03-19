
import hashlib
import urllib

from django import template


register = template.Library()


class GravatarUrlNode(template.Node):

    def __init__(self, email):
        self.email = template.Variable(email)

    def render(self, context):
        try:
            email = self.email.resolve(context)
        except template.VariableDoesNotExist:
            return ''
        # default = 'http://example.com/static/images/defaultavatar.jpg'
        default = 'mm'  # Mystery Man
        size = 24
        return '//www.gravatar.com/avatar/{}?{}'.format(
            hashlib.md5(email.lower()).hexdigest(),
            urllib.urlencode({'d': default, 's': str(size)}))


@register.tag
def gravatar_url(_parser, token):
    try:
        _tag_name, email = token.split_contents()
    except ValueError:
        raise template.TemplateSyntaxError(
            '{} tag requires a single argument'.format(
                token.contents.split()[0]))
    return GravatarUrlNode(email)
