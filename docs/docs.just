import? '../env.just'

DOCS_IMAGE_NAME := "cookiefarm-docs"
DOCS_BIN_DIR := "./"

# Build docs for prod
[group('docs')]
[working-directory('docs')]
docs-build-prod:
    @echo -e "{{ CYAN }}[*] Building docs...{{ RESET }}"
    @{{ MKDIR_CMD }} {{ DOCS_BIN_DIR }}
    @docker build -t {{ DOCS_IMAGE_NAME }} .
    @echo -e "{{ GREEN }}[+] Docs build complete!{{ RESET }}"

# Run docs for prod
[group('docs')]
[working-directory('docs')]
docs-up: docs-build-prod
    @echo -e "{{ CYAN }}[*] Docker up...{{ RESET }}"
    @{{ MKDIR_CMD }} {{ SERVER_BIN_DIR }}
    @docker run -d -p 8882:8882 --name cookiefarm-docs {{ DOCS_IMAGE_NAME }}
    @echo -e "{{ GREEN }}[+] Server build complete!{{ RESET }}"

# Stop docs for prod
[group('docs')]
[working-directory('docs')]
docs-down:
    @echo -e "{{ CYAN }}[*] Docker down...{{ RESET }}"
    @docker stop cookiefarm-docs || true
    @docker rm cookiefarm-docs || true
    @echo -e "{{ GREEN }}[+] Server stopped and removed!{{ RESET }}"

# Corresponding package.json scripts exposed as just targets

# These run the bun scripts defined in docs/package.json
[group('docs')]
[working-directory('docs')]
docs-build:
    @echo -e "{{ CYAN }}[*] Running bun run build...{{ RESET }}"
    @bun run build
    @echo -e "{{ GREEN }}[+] bun build finished{{ RESET }}"

[group('docs')]
[working-directory('docs')]
docs-dev:
    @echo -e "{{ CYAN }}[*] Running bun run dev...{{ RESET }}"
    @bun run dev

[group('docs')]
[working-directory('docs')]
docs-start:
    @echo -e "{{ CYAN }}[*] Running bun run start...{{ RESET }}"
    @bun run start

[group('docs')]
[working-directory('docs')]
docs-types-check:
    @echo -e "{{ CYAN }}[*] Running bun run types:check...{{ RESET }}"
    @bun run "types:check"

[group('docs')]
[working-directory('docs')]
docs-postinstall:
    @echo -e "{{ CYAN }}[*] Running bun run postinstall...{{ RESET }}"
    @bun run postinstall

[group('docs')]
[working-directory('docs')]
docs-lint:
    @echo -e "{{ CYAN }}[*] Running bun run lint...{{ RESET }}"
    @bun run lint
