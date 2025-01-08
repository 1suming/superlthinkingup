import {FC,useState, useEffect } from "react";
import throttle from "lodash/throttle";

import "./ScrollToTopBtn.scss"

const ScrollToTopBtn :FC= () => {
    const [isVisible, setIsVisible] = useState<boolean>(false); // 使用 TypeScript 的类型约束

    useEffect(() => {
        const handleScroll = throttle(() => {
          setIsVisible(window.scrollY > 300);
        }, 200); // 每 200 毫秒触发一次

    window.addEventListener("scroll", handleScroll);
    return () => {
        window.removeEventListener("scroll", handleScroll);
    };
    }, []);

    const scrollToTop = () => {
        window.scrollTo({
            top: 0,
            behavior: "smooth",
        });
    };

  return (
    <div>
      {isVisible && (
        <button onClick={scrollToTop} className="scroll-to-top-btn" data-tooltip="返回顶部">
          <i className="bi bi-caret-up-fill" ></i>

           
        </button>
      )}
    </div>
  );
};

export default ScrollToTopBtn;
