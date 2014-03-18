"""
Helper functions used by the Deis server.
"""
import base64
import hashlib
import random


def generate_app_name():
    """Return a randomly-generated memorable name."""
    adjectives = [
        'ablest', 'absurd', 'actual', 'allied', 'artful', 'atomic', 'august',
        'bamboo', 'benign', 'blonde', 'blurry', 'bolder', 'breezy', 'bubbly',
        'candid', 'casual', 'cheery', 'classy', 'clever', 'convex', 'cubist',
        'dainty', 'dapper', 'decent', 'deluxe', 'docile', 'dogged', 'drafty',
        'earthy', 'easier', 'edible', 'elfish', 'excess', 'exotic', 'expert',
        'fabled', 'famous', 'feline', 'finest', 'flaxen', 'folksy', 'frozen',
        'gaslit', 'gentle', 'gifted', 'ginger', 'global', 'golden', 'grassy',
        'hearty', 'hidden', 'hipper', 'honest', 'humble', 'hungry', 'hushed',
        'iambic', 'iconic', 'indoor', 'inward', 'ironic', 'island', 'italic',
        'jagged', 'jangly', 'jaunty', 'jiggly', 'jovial', 'joyful', 'junior',
        'kabuki', 'karmic', 'keener', 'kindly', 'kingly', 'klutzy', 'knotty',
        'lambda', 'leader', 'linear', 'lively', 'lonely', 'loving', 'luxury',
        'madras', 'marble', 'mellow', 'metric', 'modest', 'molten', 'mystic',
        'native', 'nearby', 'nested', 'newish', 'nickel', 'nimbus', 'nonfat',
        'oblong', 'offset', 'oldest', 'onside', 'orange', 'outlaw', 'owlish',
        'padded', 'peachy', 'pepper', 'player', 'preset', 'proper', 'pulsar',
        'quacky', 'quaint', 'quartz', 'queens', 'quinoa', 'quirky',
        'racing', 'rental', 'rising', 'rococo', 'rubber', 'rugged', 'rustic',
        'sanest', 'scenic', 'shadow', 'skiing', 'stable', 'steely', 'syrupy',
        'taller', 'tender', 'timely', 'trendy', 'triple', 'truthy', 'twenty',
        'ultima', 'unbent', 'unisex', 'united', 'upbeat', 'uphill', 'usable',
        'valued', 'vanity', 'velcro', 'velvet', 'verbal', 'violet', 'vulcan',
        'webbed', 'wicker', 'wiggly', 'wilder', 'wonder', 'wooden', 'woodsy',
        'yearly', 'yeasty', 'yeoman', 'yogurt', 'yonder', 'youthy', 'yuppie',
        'zaftig', 'zanier', 'zephyr', 'zeroed', 'zigzag', 'zipped', 'zircon',
    ]
    nouns = [
        'anaconda', 'airfield', 'aqualung', 'armchair', 'asteroid', 'autoharp',
        'babushka', 'bagpiper', 'barbecue', 'bookworm', 'bullfrog', 'buttress',
        'caffeine', 'chinbone', 'countess', 'crawfish', 'cucumber', 'cutpurse',
        'daffodil', 'darkroom', 'doghouse', 'dragster', 'drumroll', 'duckling',
        'earthman', 'eggplant', 'electron', 'elephant', 'espresso', 'eyetooth',
        'falconer', 'farmland', 'ferryman', 'fireball', 'footwear', 'frosting',
        'gadabout', 'gasworks', 'gatepost', 'gemstone', 'goldfish', 'greenery',
        'handbill', 'hardtack', 'hawthorn', 'headwind', 'henhouse', 'huntress',
        'icehouse', 'idealist', 'inchworm', 'inventor', 'insignia', 'ironwood',
        'jailbird', 'jamboree', 'jerrycan', 'jetliner', 'jokester', 'joyrider',
        'kangaroo', 'kerchief', 'keypunch', 'kingfish', 'knapsack', 'knothole',
        'ladybird', 'lakeside', 'lambskin', 'larkspur', 'lollipop', 'lungfish',
        'macaroni', 'mackinaw', 'magician', 'mainsail', 'mongoose', 'moonrise',
        'nailhead', 'nautilus', 'neckwear', 'newsreel', 'novelist', 'nuthatch',
        'occupant', 'offering', 'offshoot', 'original', 'organism', 'overalls',
        'painting', 'pamphlet', 'paneling', 'pendulum', 'playroom', 'ponytail',
        'quacking', 'quadrant', 'queendom', 'question', 'quilting', 'quotient',
        'rabbitry', 'radiator', 'renegade', 'ricochet', 'riverbed', 'rucksack',
        'sailfish', 'sandwich', 'sculptor', 'seashore', 'seedcake', 'stickpin',
        'tabletop', 'tailbone', 'teamwork', 'teaspoon', 'traverse', 'turbojet',
        'umbrella', 'underdog', 'undertow', 'unicycle', 'universe', 'uptowner',
        'vacation', 'vagabond', 'valkyrie', 'variable', 'villager', 'vineyard',
        'waggoner', 'waxworks', 'waterbed', 'wayfarer', 'whitecap', 'woodshed',
        'yachting', 'yardbird', 'yearbook', 'yearling', 'yeomanry', 'yodeling',
        'zaniness', 'zeppelin', 'ziggurat', 'zirconia', 'zoologer', 'zucchini',
    ]
    return "{}-{}".format(
        random.choice(adjectives), random.choice(nouns))


def dict_diff(dict1, dict2):
    """
    Returns the added, changed, and deleted items in dict1 compared with dict2.

    :param dict1: a python dict
    :param dict2: an earlier version of the same python dict
    :return: a new dict, with 'added', 'changed', and 'removed' items if
             any were found.

    >>> d1 = {1: 'a'}
    >>> dict_diff(d1, d1)
    {}
    >>> d2 = {1: 'a', 2: 'b'}
    >>> dict_diff(d2, d1)
    {'added': {2: 'b'}}
    >>> d3 = {2: 'B', 3: 'c'}
    >>> expected = {'added': {3: 'c'}, 'changed': {2: 'B'}, 'deleted': {1: 'a'}}
    >>> dict_diff(d3, d2) == expected
    True
    """
    diff = {}
    set1, set2 = set(dict1), set(dict2)
    # Find items that were added to dict2
    diff['added'] = {k: dict1[k] for k in (set1 - set2)}
    # Find common items whose values differ between dict1 and dict2
    diff['changed'] = {
        k: dict1[k] for k in (set1 & set2) if dict1[k] != dict2[k]
    }
    # Find items that were deleted from dict2
    diff['deleted'] = {k: dict2[k] for k in (set2 - set1)}
    return {k: diff[k] for k in diff if diff[k]}


def fingerprint(key):
    """
    Return the fingerprint for an SSH Public Key
    """
    key = base64.b64decode(key.strip().split()[1].encode('ascii'))
    fp_plain = hashlib.md5(key).hexdigest()
    return ':'.join(a + b for a, b in zip(fp_plain[::2], fp_plain[1::2]))


if __name__ == "__main__":
    import doctest
    doctest.testmod()
