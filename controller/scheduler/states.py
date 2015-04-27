import enum


class TransitionError(Exception):
    """Raised when a transition from one state to another is illegal"""

    def __init__(self, prev, next, msg):
        self.prev = prev
        self.next = next
        self.msg = msg


class JobState(enum.Enum):
    initialized = 1
    created = 2
    up = 3
    down = 4
    destroyed = 5
    crashed = 6
    error = 7
