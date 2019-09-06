_todo_completion()
{
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts=$(todo -h | grep '*' | cut -d" " -f2)

    if [ ${COMP_CWORD} == 1 ]; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    fi
}


complete -F _todo_completion todo

# WARNING: à noter que l'utilisation de la complétion automatique
# impose de fait l'utilisation de shell bash (shell dans lequel sera
# sourcé le présent fichier de configuration), car la completion
# automatique n'est pas implémentée dans le simple sh (la commande
# complete par exemple est une primitive du bash)
