 
#--where="id=36 and post_date_*>='2013-10-25'"
#mysqldump -uroot -p --default-character-set=utf8  superthinking ta_article > ta_article_backup.sql
dateStr=`date -d -1day '+%Y%m%d'`
fileName=sql${dateStr}_
bin="mysqldump superthinking  --no-create-info   -h127.0.0.1 -P3307 -uroot  -padminRoot@8888SecretPwd --opt --default-character-set=utf8mb4 --single-transaction --skip-triggers --skip-lock-tables --tables "
tables=(
ta_article
tag 
tag_rel

)
 

for tbl_name in ${tables[@]}; do
   $bin $tbl_name > /tmp/$fileName_${tbl_name}.sql 

# do something....

done
 
 
 
