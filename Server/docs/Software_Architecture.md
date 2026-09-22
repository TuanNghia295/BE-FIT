### BE FIT Backend Software Architecture
#### Clean Architecture
1. Domain - Repositories
- Repository là nơi chứa các phương thức = interface
để truy xuất dữ liệu từ cơ sở dữ liệu hoặc các nguồn dữ liệu khác.

2. Use case
- Chứa Business logic của 1 nghiệp vụ cụ thể
- Handler / Controller
        ↓
     Usecase
        ↓
   Repository
        ↓
    Database