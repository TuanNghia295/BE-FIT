# BE-FIT Feature Roadmap

## Mục tiêu

Xây dựng backend cho ứng dụng fitness cá nhân hóa: quản lý tài khoản và hồ sơ, đặt mục tiêu dinh dưỡng, tạo kế hoạch tập luyện, theo dõi buổi tập và tiến độ cơ thể.

## Hiện trạng

- Đã có đăng ký, đăng nhập JWT ES256, refresh-token rotation và logout.
- Đã có `/me` và `PUT /me` cùng auth middleware.
- Entities đã mô tả hồ sơ, mục tiêu, lịch tập, bài tập, workout session và body progress.
- Mới có test cho auth/JWT và một phần HTTP handler.
- Migrations hiện có chủ yếu chỉnh sửa/đổi tên bảng; chưa thấy migration tạo toàn bộ schema từ database rỗng.

## Roadmap

### 0. Củng cố nền tảng

- Tạo migration baseline để khởi tạo schema mới từ database rỗng.
- Rà soát các migrations hiện có, ràng buộc foreign key, unique, nullability và index.
- Kiểm tra routing: refresh/logout cần hoạt động khi access token hết hạn; không bắt buộc access-token middleware cho các endpoint này nếu refresh token là credential xác thực.
- Chuẩn hóa validation, error responses, logging và API versioning.
- Hoàn thiện test repository/use case/handler; thêm PostgreSQL integration test cho luồng token.

**Hoàn thành khi:** database thử nghiệm có thể tạo mới chỉ bằng migrations; auth routes xử lý đúng trạng thái hết hạn, revoked và lỗi database.

### 1. Hồ sơ người dùng

- Hoàn thiện `GET /me` và `PUT /me`.
- Cho phép tạo hoặc cập nhật `FitnessProfile` gắn với user đang xác thực.
- Kiểm tra DTO có đủ trường nghiệp vụ cần thiết như tuổi, giới tính, chiều cao, cân nặng và level.
- Không nhận `userId` từ body; lấy ID từ access token.
- Không trả password/hash trong response.

**Hoàn thành khi:** user chỉ đọc/cập nhật được hồ sơ của mình; test bao gồm chưa đăng nhập, user không tồn tại, profile chưa tạo và dữ liệu không hợp lệ.

### 2. Mục tiêu và dinh dưỡng

- CRUD mục tiêu cá nhân: cutting, bulking, maintain, strength.
- Lưu lịch sử mục tiêu thay vì ghi đè mục tiêu cũ.
- Tính TDEE và đề xuất calories, protein, carbs, fat dựa trên profile và mục tiêu.
- Hiển thị hoặc lưu công thức/phiên bản tính để kết quả có thể giải thích và kiểm thử.

**Hoàn thành khi:** mục tiêu thuộc đúng user, dữ liệu dinh dưỡng có validation và công thức có unit tests.

### 3. Lịch sẵn sàng tập

- CRUD `TrainingAvailability`.
- Chuẩn hóa ngày trong tuần và định dạng giờ; kiểm tra giờ kết thúc sau giờ bắt đầu.
- Cho phép cập nhật lịch mà không tạo bản ghi trùng không mong muốn.

**Hoàn thành khi:** user quản lý được các khung giờ tập và các trường hợp lịch trùng/sai được xử lý rõ ràng.

### 4. Thư viện bài tập

- API danh sách bài tập có phân trang, tìm kiếm và lọc theo nhóm cơ, độ khó, thiết bị.
- Xác định quy trình quản trị/seed dữ liệu bài tập.
- Kiểm tra quyền truy cập nếu bổ sung API tạo/sửa/xóa bài tập.

**Hoàn thành khi:** client duyệt được catalog ổn định; API không trả quá nhiều dữ liệu trong một request.

### 5. Tạo kế hoạch tập luyện

- Sinh `WorkoutPlan` theo mục tiêu, availability, số buổi/tuần và thời lượng mỗi buổi.
- Tạo schedule và danh sách exercise/set phù hợp.
- Lưu version; tạo plan mới khi user yêu cầu generate lại thay vì âm thầm sửa lịch sử.
- Bắt đầu bằng thuật toán rule-based có thể giải thích; chỉ cân nhắc AI sau khi có dữ liệu và tiêu chí đánh giá.

**Hoàn thành khi:** plan luôn thỏa availability và duration; có test cho các cấu hình đầu vào và trường hợp không thể tạo plan.

### 6. Ghi nhận buổi tập

- API bắt đầu, xem, hoàn thành hoặc bỏ qua `WorkoutSession`.
- Kiểm tra session thuộc user hiện tại và liên kết với plan/schedule hợp lệ.
- Lưu `ExercisePerformance` theo từng bài/set.
- Chống gửi lặp khi client retry request.

**Hoàn thành khi:** có thể theo dõi trọn một buổi tập và không tạo dữ liệu trùng do retry.

### 7. Theo dõi tiến độ cơ thể

- CRUD bản ghi `BodyProgress` theo thời gian.
- API lịch sử và thống kê thay đổi cân nặng.
- Chỉ user sở hữu dữ liệu mới được xem/sửa bản ghi đó.

**Hoàn thành khi:** dữ liệu có timestamp, truy vấn theo khoảng thời gian và kiểm thử quyền sở hữu.

### 8. Thông báo và subscription

- Thêm reminder tập luyện sau khi lịch tập ổn định.
- Thêm notification preferences và cơ chế gửi có retry.
- Tích hợp subscription/payment sau khi core fitness flow hoạt động; xác minh webhook và trạng thái entitlement ở server.

**Hoàn thành khi:** notification không gửi trùng; trạng thái subscription được xác minh server-side, không tin dữ liệu do client tự khai.

## Nguyên tắc xuyên suốt

- Mỗi endpoint đi qua handler → use case → repository interface → implementation.
- Mọi truy vấn dữ liệu theo user phải giới hạn theo ID lấy từ token đã xác thực.
- Mỗi feature cần validation, domain errors, tests và migration tương ứng.
- Không thêm bảng/cột bằng cách chỉ sửa entity; schema phải được cập nhật qua migration.
- Với dữ liệu nhạy cảm, chỉ trả đúng DTO cần cho client.
- Ưu tiên hoàn thiện và kiểm thử một luồng end-to-end trước khi mở rộng feature tiếp theo.

## Thứ tự ưu tiên đề xuất

1. Nền tảng schema, routing và kiểm thử.
2. Hồ sơ người dùng.
3. Mục tiêu và dinh dưỡng.
4. Availability và exercise catalog.
5. Sinh workout plan.
6. Session, performance và body progress.
7. Notifications và subscription.