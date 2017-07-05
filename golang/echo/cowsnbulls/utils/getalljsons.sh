for i in {2..12}
do
    
    curl https://www.wordgamedictionary.com/word-lists/$i-letter-words/$i-letter-words.json > ../data/$i-letterwords.json
done