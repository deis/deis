import enum


class TransitionNotAllowed(Exception):
    """Raised when a transition from one state to another is illegal"""


class JobState(enum.Enum):
    initialized = 1
    created = 2
    up = 3
    down = 4
    destroyed = 5
    crashed = 6
    error = 7
