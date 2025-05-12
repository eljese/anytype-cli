# Anytype CLI

## Run headless MW server through CLI

Expected repository structure:

```
parent-directory/
├── anytype-heart/
└── anytype-cli/
```

1. **In `anytype-heart` (go-4643-headless-client-anytype-cli) directory:**

```bash
make install-dev-cli
```

2. **In `anytype-cli` directory:**

```bash
make build
```

3. **Start the daemon:**

- To run in the foreground:

```bash
./dist/anytype-cli daemon
```

- To run in the background:

```bash
./dist/anytype-cli daemon &
```

## Auto-approve members in a space

```bash
./dist/anytype-cli space autoapprove --role "Editor" --space "<SpaceId>"
```
