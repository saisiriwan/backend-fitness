# POstgreSQL config
/*หน้าที่: นี่คือ "คลังวัตถุดิบ" 
ไฟล์สำคัญ:
docker-compose.yml : เป็นพิมพ์เขียวสำหรับรัน PostgreSQL Database (db) และ PgAdmin (เว็บสำหรับดูข้อมูลใน DB) 
docker/Dockerfile : ใช้สำหรับสร้าง Image ของ PostgreSQL
docker/init.sql (ไฟล์นี้ถูกอ้างอิง): นี่คือไฟล์ที่เราใช้สร้างตาราง users, clients, workout_sessions ฯลฯ ตามที่เราออกแบบไว้
backup/: นี่คือ Service เสริมสำหรับ Backup ฐานข้อมูล*/
