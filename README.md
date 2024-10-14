# BackEnd-Gift


# Steps
1. Create the PostgreSQL database with a name that could be linked to the Golang API using a PSQL management tool (I used psql Shell and pgAdmin).
   
2. Install dependencies related to the framework and tools/libraries used (I used Gin framework with GORM for Object-Relational Mapping, Redis for Caching, and bluemonday for XSS attack prevention).

3. Creating project folders with main folder for running the main function, while hiding the details for extra privacy inside the module folder where parts of the project into one folder each such as for the article-related tables’ models (named articles), another for controllers, data, middlewares, and the router standing on inside the root module folder to allow easier access to the other components.  

4. Import the GORM tool and PostgreSQL driver in the database.go file (make sure both have already been installed beforehand via Terminal in step 2) which was optionally separated into the data folder.

5. Create the database schema for the Articles and ArticleCategories tables, one each in separate files, with one-to-many relationship from ArticleCategories to Articles. They must include JSON and GORM notations related to each data's attributes later relevant for data insertion and validation (the Category GORM in Articles model has an "OnDelete:CASCADE;" constraint in particular so that it would also be deleted alongside the CategoryID from the ArticleCategories table when a Delete request is sent).

6. Create controllers containing each of the tables' CRUD logic functions using GORM as a connection to the database. The database schema must be imported in order to connect the CRUD functions to each data properly, while Create (Add) and Update functions should include bluemonday sanitization to prevent XSS attacks, then the Read function includes pagination to divide 10 data into a single page accessible via a new query which allows users to see which page of the data are they seeing.

7. Create middlewares for handling ID validation (idmiddleware.go to make sure that the automatically set and incremented ID does not bring in a non-numeric, zero, or negative value) and rate-limiting (ratelimitmiddleware.go to prevent one or some users' domination over the server which results in the server crashing/Denial of Service). 

8. Create router to connect the controllers and middlewares with the endpoints, grouped according to the destination table (articles/categories) and whether one (:id) or all of the data are selected. Add functions with Redis caching are included in this file to improve database load time.

9. Testing and Documenting the endpoints via Postman. API Documentation available via this link: https://www.postman.com/spacecraft-cosmologist-96014201/workspace/my-workspace/collection/38938614-85716a75-70c0-48fc-9301-7019800381b8?action=share&creator=38938614
   
# Consideration Answers
1. What if there are thousands of articles in the database?
   Create pagination which prevents crashing/data overload through not showing all data at a single time, instead dividing them to 10 per page.

   Codes in GetArticle function in articles_controllers.go related to this consideration:
    *pageStr := c.Query("page")
	  pageSizeStr := c.Query("page_size")
	  page, _ := strconv.Atoi(pageStr)
	  if page < 1 {
		  page = 1
	  }
	  pageSize, _ := strconv.Atoi(pageSizeStr)
	  if pageSize < 1 {
      pageSize = 10
	  }
	  offset := (page - 1) * pageSize
	  var articles []articles.Articles
	  if err := ac.DB.Preload("Category").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch!"})
		  return
	  }*

   Codes in GetCategories function in article_categories_controllers.go related to this consideration:
    *pageStr := c.Query("page")
	  pageSizeStr := c.Query("page_size")
	  page, _ := strconv.Atoi(pageStr)
	  if page < 1 {
		  page = 1
	  }
	  pageSize, _ := strconv.Atoi(pageSizeStr)
	  if pageSize < 1 {
		  pageSize = 10
	  }
	  offset := (page - 1) * pageSize
	  var categories []articles.ArticleCategories
	  if err := ac.DB.Preload("Category").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&categories).Error; err != nil {
		  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch!"})
		  return
    }*

   2. What if many users are accessing your API at same time?
      The ratelimitmiddleware.go works exactly to limit request from the same user detected by their IP (due to the small test API server, it is set to 10 requests per user before they receive a HTTP 429 Too Many Requests error and get         blocked for a minute before their request count resets) to prevent server overload so that not only a single user or bot dominates the server, but also having other users possibly accessing the API at the same time without any            techinical difficulties.

      Codes in ratelimitmiddleware.go related to this consideration:
      *var (
        	rateLimit    = 5
        	rateReset    = time.Minute
        	userRequests = make(map[string]int)
        	mu           sync.Mutex
        )
        func RateLimitMiddleware() gin.HandlerFunc {
        	return func(c gin.Context) {
        		userIP := c.ClientIP()
        		mu.Lock()
        		if count, exists := userRequests[userIP]; exists {
        			if count >= rateLimit {
        				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded, try again later"})
        				mu.Unlock()
        				c.Abort()
        				return
        			}
        			userRequests[userIP]++
        		} else {
        			userRequests[userIP] = 1
        		}
        		mu.Unlock()
        		c.Next()
        		go func() {
        			time.Sleep(rateReset)
        			mu.Lock()
        			userRequests[userIP]--
        			if userRequests[userIP] <= 0 {
        				delete(userRequests, userIP)
        			}
        			mu.Unlock()
        		}()
        	}
        }*

      3. What if users perform stored xss and how to prevent it?
         The sanitized user input offerred by the bluemonday library ensures that none of the malicious threats such as malware or data theft can be sent via XSS attacks and injections. The UGCPolicy in particular checks and sanitizes             any user-inputted HTML which means that no harmful HTML links could be stored and displayed in the API as well.

         Codes in AddArticle and UpdateArticle functions in articles_controller.go related to this consideration:
         *p := bluemonday.UGCPolicy()
	        article.Title = p.Sanitize(article.Title)
	        article.Content = p.Sanitize(article.Content)
	        article.Thumbnail = p.Sanitize(article.Thumbnail)
	        article.Slug = p.Sanitize(article.Slug)*
         
         Codes in AddCategory and UpdateCategory functions in article_categories_controller.go related to this consideration:
         *p := bluemonday.UGCPolicy()
	        categories.Name = p.Sanitize(categories.Name)
	        categories.Slug = p.Sanitize(categories.Slug)*
