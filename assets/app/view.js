document.addEventListener('DOMContentLoaded', () => {
    const deleteButton = document.querySelector('button[data-action="delete-post"]');
    if (!deleteButton) return;

    deleteButton.addEventListener('click', () => {
        const postID = deleteButton.dataset.postId;

        if (confirm("Are you sure you want to delete this post? This action cannot be undone.")) {
            fetch(`/view?post_id=${postID}`, {
                method: 'DELETE',
                headers: {
                    'csrf': document.getElementById('csrf').value,
                },
            })
                .then(response => {
                    if (response.ok) {
                        alert("Post deleted successfully.");
                        window.location.href = '/home';
                    } else {
                        alert("Failed to delete the post.");
                    }
                })
                .catch(error => {
                    console.error("Error:", error);
                    alert("An error occurred while deleting the post.");
                });
        }
    });

    // Comment delete logic
    document.querySelectorAll('button[data-action="delete-comment"]').forEach(btn => {
        btn.addEventListener('click', () => {
            const commentID = btn.dataset.commentId;
            if (confirm("Are you sure you want to delete this comment? This action cannot be undone.")) {
                fetch(`/view?comment_id=${commentID}`, {
                    method: 'DELETE',
                    headers: {
                        'csrf': document.getElementById('csrf').value,
                    },
                })
                    .then(response => {
                        if (response.ok) {
                            // Remove the comment from the UI, or reload
                            btn.closest('.items-center.w-full.justify-center.flex').remove();
                        } else {
                            alert("Failed to delete the comment.");
                        }
                    })
                    .catch(error => {
                        console.error("Error:", error);
                        alert("An error occurred while deleting the comment.");
                    });
            }
        });
    });

});

document.addEventListener('DOMContentLoaded', () => {
    const csrfToken = document.getElementById('csrf')?.value || '';

    document.querySelectorAll('button[data-action="edit-comment"]').forEach(button => {
        button.addEventListener('click', function() {
            const commentID = this.dataset.commentId;
            const container = this.closest('.comment-container');
            const contentDD = container.querySelector('.comment-content');
            if (container.querySelector('textarea')) return;

            const originalContent = contentDD.textContent.trim();

            // Replace <dd> with textarea and buttons (below)
            const editFrame = document.createElement('div');
            editFrame.className = "w-full flex flex-col mt-2"; // Ensures vertical layout

            const textarea = document.createElement('textarea');
            textarea.className = "w-full px-3 py-2 border border-gray-300 rounded-md resize-none text-sm focus:outline-none focus:ring-2 focus:ring-blue-400 mb-2";
            textarea.rows = 2;
            textarea.value = originalContent;

            // Save and Cancel
            const buttonsDiv = document.createElement('div');
            buttonsDiv.className = "flex flex-row gap-x-2"; // Button group horizontally

            const saveBtn = document.createElement('button');
            saveBtn.textContent = "Save";
            saveBtn.type = "button";
            saveBtn.className = "px-2 py-1 bg-green-100 ring ring-green-600 ring-inset rounded-lg hover:bg-green-400 transition";

            const cancelBtn = document.createElement('button');
            cancelBtn.textContent = "Cancel";
            cancelBtn.type = "button";
            cancelBtn.className = "px-2 py-1 bg-blue-100 ring ring-blue-600 ring-inset rounded-lg hover:bg-blue-400 transition";

            buttonsDiv.appendChild(saveBtn);
            buttonsDiv.appendChild(cancelBtn);

            editFrame.appendChild(textarea);
            editFrame.appendChild(buttonsDiv);

            // Insert the edit frame BEFORE the contentDD so it's in the same visual block
            contentDD.style.display = 'none';
            contentDD.parentNode.insertBefore(editFrame, contentDD);

            // Hide Edit button while editing
            this.style.display = 'none';

            textarea.focus();

            cancelBtn.addEventListener('click', () => {
                editFrame.remove();
                contentDD.style.display = '';
                button.style.display = '';
            });

            saveBtn.addEventListener('click', () => {
                const newContent = textarea.value.trim();
                const formData = new FormData();
                formData.append('content', newContent);
                formData.append('csrf', csrfToken);

                fetch(`/view?comment_id=${commentID}`, {
                    method: 'PUT',
                    body: formData
                })
                    .then(res => {
                        if(res.ok) {
                            contentDD.textContent = newContent;
                            editFrame.remove();
                            contentDD.style.display = '';
                            button.style.display = '';
                        } else {
                            return res.text().then(txt => { throw new Error(txt) });
                        }
                    }).catch(err => {
                    alert("Failed to update comment:\n" + err.message);
                });
            });
        });
    });

    const editButton = document.querySelector('button[data-action="edit-post"]');
    if (editButton) {
        editButton.addEventListener('click', () => {
            const postID = editButton.dataset.postId;
            window.location.href = `/edit?id=${encodeURIComponent(postID)}`;
        });
    }

});

document.addEventListener('DOMContentLoaded', () => {
    // Initialize button styles based on data-vote-state
    document.querySelectorAll('[data-vote-state]').forEach(button => {
        const state = button.dataset.voteState;
        if (state === 'upvote' && button.id.startsWith('upvote-')) {
            button.classList.add('bg-blue-300');
        } else if (state === 'downvote' && button.id.startsWith('downvote-')) {
            button.classList.add('bg-red-200');
        }
    });

    // Function to send a vote and update the UI dynamically
    const sendVote = (id, isPost, voteType) => {
        const formData = new FormData();
        const csrf = document.getElementById('csrf');
        formData.append(isPost ? 'post_id' : 'comment_id', id);
        formData.append('vote_type', voteType);
        formData.append('csrf', csrf.value);

        const postAuthor = document.getElementById('post-author-id');
        const actionAuthor = document.getElementById('action-author-id');
        if (postAuthor) {
            formData.append('post_author_id', postAuthor.value);
        }
        if (actionAuthor) {
            formData.append('action_author_id', actionAuthor.value);
        }


        fetch('/vote', {
            method: 'POST',
            body: formData,
        }).then(response => {
            if (response.ok) {
                response.json().then(data => {
                    // Update the vote counts dynamically
                    const upvoteButton = document.getElementById(`upvote-${isPost ? 'post' : 'comment'}-${id}`);
                    const downvoteButton = document.getElementById(`downvote-${isPost ? 'post' : 'comment'}-${id}`);
    
                    // Update the vote counts inside the <span> elements
                    upvoteButton.querySelector('.vote-count').textContent = data.upvotes;
                    downvoteButton.querySelector('.vote-count').textContent = data.downvotes;
    
                    // Reset button styles
                    upvoteButton.classList.remove('bg-blue-300');
                    downvoteButton.classList.remove('bg-red-200');
    
                    // Apply the new style based on the vote state
                    if (data.user_vote === 'upvote') {
                        upvoteButton.classList.add('bg-blue-300');
                    } else if (data.user_vote === 'downvote') {
                        downvoteButton.classList.add('bg-red-200');
                    }
                });
            } else {
                console.error('Failed to register vote');
                alert('Failed to register your vote. Please try again.');
            }
        }).catch(error => {
            console.error('Error:', error);
            alert('An error occurred while processing your vote.');
        });
    };

    // Add event listeners for upvote buttons
    document.querySelectorAll('[id^="upvote-"]').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.id.split('-')[2];
            const isPost = button.id.includes('post');
            sendVote(id, isPost, 'upvote');
        });
    });

    // Add event listeners for downvote buttons
    document.querySelectorAll('[id^="downvote-"]').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.id.split('-')[2];
            const isPost = button.id.includes('post');
            sendVote(id, isPost, 'downvote');
        });
    });
});