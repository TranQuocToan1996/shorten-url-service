TODO: Lập list ưu tiên các tiêu chí
Để đánh giá **một hệ thống URL rút gọn** hoặc code cho URL rút gọn là tốt, bạn có thể xem xét các tiêu chí sau:

---

1. **Deterministic / Consistency**

   * Cùng một URL phải luôn ra cùng một short code nếu muốn tránh trùng lặp trong DB.
   * Hệ thống nên kiểm tra tồn tại trước khi tạo mới.

2. **Độ ngắn và hiệu quả**

   * Short code càng ngắn càng dễ chia sẻ.
   * Base62 hoặc Base64 URL-safe là phổ biến.
   * Tránh quá ngắn → dễ đoán / brute-force.

3. **Độ bảo mật**

   * Không để short code lộ thông tin nội bộ (ID, số thứ tự).
   * Sử dụng HMAC hoặc salt để tránh đoán URL.
   * Kiểm soát URL nguy hiểm (malware/phishing).

4. **Hiệu năng / Scalability**

   * Hệ thống phải xử lý hàng triệu URL nhanh chóng.
   * Tránh truy vấn DB quá nhiều → có thể cache mapping.
   * Có thể phân tán / shard DB nếu lượng URL lớn.

5. **Collision handling**

   * Phải có cách xử lý khi hash/short code trùng nhau (có thể append ký tự, rehash…).

6. **Quản lý metadata**

   * Dễ tracking số click, nguồn truy cập.
   * Có thể thêm expiry time, owner info, hoặc analytics.

7. **Khả năng mở rộng / maintainability**

   * Code phải dễ maintain, dễ mở rộng tính năng mới.
   * Cấu trúc rõ ràng, modular, testable.

8. **Kiểm soát lỗi và logging**

   * Xử lý URL không hợp lệ, lỗi encode/decode.
   * Log đầy đủ nhưng không lộ thông tin nhạy cảm.

9. **User experience**

   * Redirect nhanh, không delay.
   * Có trang preview nếu cần.

10. **Thống kê và báo cáo**

    * Cho phép xem URL nào được click nhiều, từ đâu, thời gian…
    * Hỗ trợ phân tích mà không làm chậm redirect.

---


TODO: Các thuật toán ưu nhược

TODO: Tổng quan về lựa chọn queue


TODO: Provide detailed instructions on how to run your assignment in a separate markdown file.


TODO: Provide tests for both endpoints (and any other tests you may want to write).

TODO: You need to think through potential attack vectors on the application, and document them in the README.

TODO: Document explain how to scale(Add infra part), how solve collision problem, how solve concurrency problem

TODO: Golang best practices explain