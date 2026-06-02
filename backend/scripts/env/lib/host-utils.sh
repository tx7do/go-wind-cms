#!/usr/bin/env bash
################################################################################
## Hosts utilities (Unix)
##
## Manage /etc/hosts records by adding/removing IP-domain mappings.
## Usage:
##   source "${SCRIPT_DIR}/lib/host-utils.sh"
##
## Encoding: UTF-8 (LF) | Compatible: bash 4.0+
################################################################################

# Prevent duplicate loading
if [ "${_HOST_UTILS_LOADED:-0}" -eq 1 ]; then
  return 0 2>/dev/null || true
fi
_HOST_UTILS_LOADED=1

# Load common utils if not already loaded
if ! command -v log >/dev/null 2>&1; then
  _HOST_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  # shellcheck source=common-utils.sh
  source "${_HOST_LIB_DIR}/common-utils.sh"
fi

_hosts_can_write() {
  [ -w "/etc/hosts" ] || [ "$EUID" -eq 0 ] || [ -n "${SUDO:-}" ]
}

_flush_dns_cache() {
  if [[ "$OSTYPE" == "darwin"* ]]; then
    ${SUDO} dscacheutil -flushcache >/dev/null 2>&1 || true
    ${SUDO} killall -HUP mDNSResponder >/dev/null 2>&1 || true
    return 0
  fi

  if command -v resolvectl >/dev/null 2>&1; then
    ${SUDO} resolvectl flush-caches >/dev/null 2>&1 || true
  elif command -v systemd-resolve >/dev/null 2>&1; then
    ${SUDO} systemd-resolve --flush-caches >/dev/null 2>&1 || true
  elif command -v service >/dev/null 2>&1; then
    ${SUDO} service nscd restart >/dev/null 2>&1 || true
    ${SUDO} service dnsmasq restart >/dev/null 2>&1 || true
  fi
}

# edit_hosts <ip> <domain> [Add|Remove]
edit_hosts() {
  local ip=${1:-}
  local domain=${2:-}
  local operate=${3:-Add}
  local hosts_file="/etc/hosts"

  if [ -z "$ip" ] || [ -z "$domain" ]; then
    error "Usage: edit_hosts <ip> <domain> [Add|Remove]"
    return 1
  fi

  if ! _hosts_can_write; then
    error "Please run as root (or with sudo privileges) to edit /etc/hosts"
    return 1
  fi

  case "$operate" in
    Add|add)
      if awk -v ip="$ip" -v domain="$domain" '
        /^[[:space:]]*#/ { next }
        NF >= 2 && $1 == ip && $2 == domain { found = 1 }
        END { exit(found ? 0 : 1) }
      ' "$hosts_file"; then
        log "Record already exists, skip: $ip $domain"
        return 0
      fi

      printf '\n%s %s\n' "$ip" "$domain" | ${SUDO} tee -a "$hosts_file" >/dev/null
      success "Added: $ip $domain"
      ;;
    Remove|remove)
      local tmp_file
      tmp_file=$(mktemp)

      if awk -v ip="$ip" -v domain="$domain" '
        {
          if ($0 ~ /^[[:space:]]*#/) { print; next }
          if (NF >= 2 && $1 == ip && $2 == domain) { removed = 1; next }
          print
        }
        END { exit(removed ? 0 : 2) }
      ' "$hosts_file" > "$tmp_file"; then
        ${SUDO} tee "$hosts_file" < "$tmp_file" >/dev/null
        rm -f "$tmp_file"
        success "Removed: $ip $domain"
      else
        local status=$?
        rm -f "$tmp_file"
        if [ "$status" -eq 2 ]; then
          warn "Record does not exist, skip: $ip $domain"
          return 0
        fi
        error "Failed to remove record: $ip $domain"
        return 1
      fi
      ;;
    *)
      error "Invalid operate: $operate (expected Add or Remove)"
      return 1
      ;;
  esac

  _flush_dns_cache
  log "DNS cache refresh attempted"
  return 0
}

# initialize_hosts <ip> <domain_suffix> <service1> [service2 ...]
initialize_hosts() {
  local ip=${1:-127.0.0.1}
  local domain_suffix=${2:-.local}
  shift 2 || true

  if [ "$#" -eq 0 ]; then
    error "Usage: initialize_hosts <ip> <domain_suffix> <service1> [service2 ...]"
    return 1
  fi

  log "========== Initialize hosts records =========="

  local success_count=0
  local fail_count=0
  local service domain

  for service in "$@"; do
    domain="${service}${domain_suffix}"
    if edit_hosts "$ip" "$domain" "Add"; then
      success_count=$((success_count + 1))
    else
      fail_count=$((fail_count + 1))
    fi
  done

  log ""
  log "Done: success ${success_count}, failed ${fail_count}"

  [ "$fail_count" -eq 0 ]
}

