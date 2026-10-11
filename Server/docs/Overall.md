Các bước cấu hình, phát triển 1 sản phẩm ứng dụng mobile end to end
BE FIT

1.	Mục tiêu
-	Xây dựng 1 ứng dụng mobile giúp người dùng tập luyện, quản lý, lên lịch tập phù hợp => cá nhân hóa từng người
2.	Đối tượng
-	Người mới tập ( < 12 tháng )
-	Người đã có kinh nghiệm tập luyện ( > 12 tháng )
-	Người đã tập lâu năm ( > 3 năm )
3.	Công nghệ thiết kế
-	Backend: Golang, Postgre, Redis, Docker, CI/CD git actions
-	Frontend: React Native, React, SWR, Zustand
-	Payment: RevenueCat
-	Testing: Unit Test, Integration Test, Performance test, stress test
-	System Design: Rate limiter, Notification system, Payment,

 Các bước thực hiện đối với Backend
-	Tạo database postgre với docker trước
-	Xây dựng Software Architecture ( structure code ): Ở đây chọn Clean Architecture
-	Setup biến môi trường, git ignore
-	Kết nối DB với backend
-	Cấu hình HTTP Server
-	Thiết kế Database ( ERD )
-	Bắt  đầu thiết kế các API.
Các bước thực hiện đối với FrontEnd
-	Tạo dự án
-	Xây dựng Software Architecture: Ở đây chọn Clean Architecture
-	UI Design method: Atomic Design
-	Kết nối với Backend

 Kiến thức học được qua dự án
1.	Docker
-	Docker Volume: Tồn tại như 1 không gian lưu trữ dữ liệu độc lập. Xóa container nhưng dữ liệu trong volume vẫn còn cho đến khi volume đó được xóa
2.	Architecture
2.1	Clean architecture:
o	Chia backend ra thành các layer chính:
	Domain
•	Chứa entity và business rules
•	Là core hệ thống
•	Không phụ thuộc framework/database
	Use case
•	Chứa business logic
•	Điều phối các nghiệp vụ của hệ thống
•	Gọi Repository thông qua các interface
	Delivery
•	Nhận và trả request/response (HTTP, gRPC Handler)
	Repository / Infrastructure
•	Làm việc với Postgre, Redis,…
•	Implement các interface mà use case cần

 Flow hoạt động
1.	Người dùng chọn mục đích tập luyện. Ví dụ: goal: Cutting => TargetWeight => Tính toán bằng công thức tính TDEE => Đưa ra số Macro hợp lý (Protein,Carb,Fat)
2.	Người dùng chọn số buổi có thể tập trong tuần + Thời gian rảnh để tập => Hệ thống tự động chọn plan tập hợp lý 

