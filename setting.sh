export BCD_DIR=~/.bettercd
export BCD_FILTER_TYPE=peco
export BCD_LOG_DIR=$BCD_DIR/directory.log
export BCD_LOG_MAX=200

function bettercd(){
   local dir="$( $BCD_DIR/bettercd -i $BCD_LOG_MAX -l $BCD_LOG_DIR $@)"

   if [ "$dir" != "" ]; then
      local num="$(echo $dir | wc -l)"
      if [ $num -eq 1 ]; then
         builtin cd "$(echo $dir)"
         printf "jumped to "
         pwd
      else
         local cdir="$(pwd)"
         builtin cd "$(echo $dir | $BCD_FILTER_TYPE)"
         if [ "$cdir" != "$(pwd)" ]; then
           $BCD_DIR/bettercd -w $(pwd)
           printf "jumped to "
           pwd
         fi
      fi
   else
       cd $@
   fi
}

alias cd=bettercd
