SHELL_SCRIPTS = $(wildcard util/*) \
								$(shell find "." -name '*.sh')

EXCLUDE = util/deis-contributions-since

test-style:
	shellcheck $(filter-out $(EXCLUDE), $(SHELL_SCRIPTS))
