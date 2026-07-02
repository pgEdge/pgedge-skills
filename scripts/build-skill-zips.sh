#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SKILLS_DIR="$REPO_ROOT/skills"
OUT_DIR="${1:-$REPO_ROOT/dist}"

mkdir -p "$OUT_DIR"

for skill_dir in "$SKILLS_DIR"/*/; do
  skill=$(basename "$skill_dir")
  [ "$skill" = "shared" ] && continue          # never a standalone zip

  zip_path="$OUT_DIR/$skill.zip"
  rm -f "$zip_path"

  # Skill's own files (exclude macOS cruft)
  (cd "$skill_dir" && \
    find . -type f -not -name '.DS_Store' | zip "$zip_path" -@)

  # Bundle shared/ as a sibling dir IF it exists (future-proof)
  if [ -d "$SKILLS_DIR/shared" ]; then
    (cd "$SKILLS_DIR" && \
      find shared -type f -not -name '.DS_Store' | zip "$zip_path" -@)
  fi

  echo "  $skill.zip ($(du -h "$zip_path" | cut -f1 | xargs))"
done

echo ""
echo "Built $(ls "$OUT_DIR"/*.zip 2>/dev/null | wc -l | xargs) zips in $OUT_DIR"
