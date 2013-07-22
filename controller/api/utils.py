"""
Helper functions used by the Deis server.
"""

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
        'iambic', 'indoor', 'inward', 'italic',
        'jagged', 'jangly', 'jovial', 'junior',
        'kabuki', 'karmic', 'kindly', 'knotty',
        'leaner', 'lonely', 'loving', 'luxury',
        'madras', 'marble', 'mellow', 'molten',
        'native', 'nested', 'newish', 'nickel',
        'oblong', 'oldest', 'orange', 'owlish',
        'padded', 'peachy', 'pepper', 'proper',
        'quaint', 'quartz',
        'racing', 'rising', 'rubber', 'rugged',
        'sanest', 'scenic', 'shadow', 'skiing',
        'taller', 'tender', 'timely', 'truthy',
        'unbent', 'unisex', 'uphill', 'usable',
        'valued', 'velvet', 'violet', 'vulcan',
        'webbed', 'wicker', 'wilder', 'woodsy',
        'yearly', 'yogurt', 'yonder',
        'zeroed', 'zigzag', 'zipped', 'zonked',
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
    return '{}-{}'.format(
        random.choice(adjectives), random.choice(nouns))
