package models

import "testing"

const (
	str255char = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	str256char = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	str64char = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	str63char = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

func Test_NewPost_EmptyMessageReturnsError(t *testing.T) {
	inputs := []string{"", " ", "   "}

	for i := range inputs {
		input := inputs[i]

		post := NewPost(input, "username", true)
		errs := post.ValidatePost()

		if len(errs) == 0 {
			t.Error("post.ValidatePost(); expected errors, returned none")
		} else if _, ok := errs["message"]; ok == false {
			t.Error(`post.ValidatePost(); return errors, had no "message" key`)
			for k, v := range errs {
				t.Logf(`found errs["%s"] = "%s"`, k, v)
			}
		}
	}
}

func Test_NewPost_MessageLen1_NoError(t *testing.T) {
	input := "a"

	post := NewPost(input, "username", true)
	errs := post.ValidatePost()

	if len(errs) != 0 {
		t.Error("post.ValidatePost(); unexpected errors")
		for k, v := range errs {
			t.Logf(`found errs["%s"] = "%s"`, k, v)
		}
	} 
}

func Test_NewPost_MessageLen255_NoError(t *testing.T) {
	input := str255char

	post := NewPost(input, "username", true)
	errs := post.ValidatePost()

	if len(errs) != 0 {
		t.Error("post.ValidatePost(); unexpected errors")
		for k, v := range errs {
			t.Logf(`found errs["%s"] = "%s"`, k, v)
		}
	} 
}

func Test_NewPost_MessageLen256_Errors(t *testing.T) {
	input := str256char

	post := NewPost(input, "username", true)
	errs := post.ValidatePost()

	if len(errs) == 0 {
		t.Error("post.ValidatePost(); expected errors, got none")
	} else if _, ok := errs["message"]; ok == false {
		t.Error(`post.ValidatePost(); return errors, had no "message" key`)
		for k, v := range errs {
			t.Logf(`found errs["%s"] = "%s"`, k, v)
		}
	}
}

func Test_NewPost_EmptyUsernameReturnsError(t *testing.T) {
	inputs := []string{"", " ", "   "}

	for i := range inputs {
		input := inputs[i]

		post := NewPost("message", input, true)
		errs := post.ValidatePost()

		if len(errs) == 0 {
			t.Error("post.ValidatePost(); expected errors, returned none")
		} else if _, ok := errs["message"]; ok == false {
			t.Error(`post.ValidatePost(); return errors, had no "username" key`)
			for k, v := range errs {
				t.Logf(`found errs["%s"] = "%s"`, k, v)
			}
		}
	}
}

func Test_NewPost_UsernameLen1_NoError(t *testing.T) {
	input := "a"

	post := NewPost("message", input, true)
	errs := post.ValidatePost()

	if len(errs) != 0 {
		t.Error("post.ValidatePost(); unexpected errors")
		for k, v := range errs {
			t.Logf(`found errs["%s"] = "%s"`, k, v)
		}
	} 
}

func Test_NewPost_UsernameLen63_NoError(t *testing.T) {
	input := str63char

	post := NewPost("message", input, true)
	errs := post.ValidatePost()

	if len(errs) != 0 {
		t.Error("post.ValidatePost(); unexpected errors")
		for k, v := range errs {
			t.Logf(`found errs["%s"] = "%s"`, k, v)
		}
	} 
}

func Test_NewPost_MessageLen64_Errors(t *testing.T) {
	input := str64char

	post := NewPost("message", input, true)
	errs := post.ValidatePost()

	if len(errs) == 0 {
		t.Error("post.ValidatePost(); expected errors, got none")
	} else if _, ok := errs["message"]; ok == false {
		t.Error(`post.ValidatePost(); return errors, had no "username" key`)
		for k, v := range errs {
			t.Logf(`found errs["%s"] = "%s"`, k, v)
		}
	}
}
