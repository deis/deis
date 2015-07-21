function clean-vagrant {
  "${THIS_DIR}/halt-all-vagrants.sh"
  vagrant destroy --force
}
