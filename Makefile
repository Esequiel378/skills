.PHONY: install list

install: ## Symlink all skills into ~/.claude/skills (install/update)
	@./scripts/link-skills.sh

list: ## Print each skill's name + description
	@./scripts/list-skills.sh
