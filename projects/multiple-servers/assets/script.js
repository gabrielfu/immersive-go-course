function fetchImages(development) {
    if (development) {
      return Promise.resolve([
        {
          title: "Sunset",
          alt_text: "Clouds at sunset",
          url: "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
        },
        {
          title: "Mountain",
          alt_text: "A mountain at sunset",
          url: "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
        },
      ]);
    }
    return fetch("http://localhost:8081/images.json").then((_) => _.json());
  }
  
  function timeout(t, v) {
    return new Promise((res) => {
      setTimeout(() => res(v), t);
    });
  }
  
  const gallery$ = document.querySelector(".gallery");
  
  fetchImages(false).then(
    (images) => {
      gallery$.textContent = images.length ? "" : "No images available.";
  
      images.forEach((img) => {
        const imgElem$ = document.createElement("img");
        imgElem$.src = img.url;
        imgElem$.alt = img.alt_text;
        const titleElem$ = document.createElement("h3");
        titleElem$.textContent = img.title;
        const wrapperElem$ = document.createElement("div");
        wrapperElem$.classList.add("gallery-image");
        wrapperElem$.appendChild(titleElem$);
        wrapperElem$.appendChild(imgElem$);
        gallery$.appendChild(wrapperElem$);
      });
    },
    () => {
      gallery$.textContent = "Something went wrong.";
    }
  );